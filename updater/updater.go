package updater

import (
	"charles_communicator"
	"common"
	"crypto/sha256"
	"device_info"
	"encoding/hex"
	"errors"
	"fmt"
	"gablogger"
	"initializer"
	"io"
	"log"
	"os"
	"os/exec"
	"peripherals"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/sftp"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/ssh"
)

var Logger = gablogger.Logger()
var Handler *Updater

type SFTPConfig struct {
	host     string
	port     string
	username string
	password string
}

type Updater struct {
	localPath      string
	remotePath     string
	sftp           *SFTPConfig
	Hlk7628Topic   string
	Stm32Topic     string
	stm32Version   string
	hlk7628Version string
}

func InitUpdater() {
	Handler = &Updater{}
	Handler.Hlk7628Topic = fmt.Sprintf("environments/%s/hlk7628/version", common.ENVIRONMENT)
	Handler.Stm32Topic = fmt.Sprintf("environments/%s/stm32/version", common.ENVIRONMENT)
	Handler.sftp = newSFTPConfig(common.SFTP_SERVER, common.SFTP_PORT, common.SFTP_USER, common.SFTP_PASS)

	OSVersion, err := device_info.GetOSVersion()
	if err != nil {
		return
	}
	Logger.Infoln("Current OSVersion " + OSVersion)

	STM32Version, err := peripherals.GetFirmwareVersion()
	if err != nil {
		return
	}
	Logger.Infoln("Current STM32Version " + STM32Version)

	Handler.hlk7628Version = OSVersion
	Handler.stm32Version = STM32Version
	Handler.localPath = "/tmp"
	Handler.remotePath = "Files"
}

func UpdaterHlk7628Callback(client mqtt.Client, message mqtt.Message) {
	version := gjson.Get(string(message.Payload()), "version").String()
	sha256sum := gjson.Get(string(message.Payload()), "sha256sum").String()
	Logger.Debugln("Remote OSVersion " + version)

	if strings.Compare(version, Handler.hlk7628Version) != 0 {
		Logger.Infoln("Updating hlk7628 version from " + Handler.hlk7628Version + " to " + version)
		remotePath := Handler.remotePath + "/hlk7628/" + common.ENVIRONMENT + "/" + version + "/charlinhos-sysupgrade.bin"
		localPath := Handler.localPath + "/charlinhos-sysupgrade.bin"
		updateDevice(remotePath, localPath, sha256sum, applyHlk7628Update)
	} else {
		callHLK7628Events("updated")
	}
}

func UpdaterStm32Callback(client mqtt.Client, message mqtt.Message) {
	version := gjson.Get(string(message.Payload()), "version").String()
	sha256sum := gjson.Get(string(message.Payload()), "sha256sum").String()
	Logger.Debugln("Remote STM32Version " + version)

	if strings.Compare(version, Handler.stm32Version) != 0 {
		Logger.Infoln("Updating stm32 version from " + Handler.stm32Version + " to " + version)
		remotePath := Handler.remotePath + "/stm32/" + common.ENVIRONMENT + "/" + version + "/firmware.bin"
		localPath := Handler.localPath + "/firmware.bin"
		charles_communicator.ClosePort()
		updateDevice(remotePath, localPath, sha256sum, applyStm32Update)
		charles_communicator.OpenPort()
		Logger.Infoln("STM Updated")
		Handler.stm32Version = version
	} else {
		callSTM32Events("updated")
	}
}

// updateDevice downloads a file, verifies it, and updates a device using the provided update function.
//
// Parameters:
//
//	remotePath: The remote path of the file to be downloaded.
//	localPath: The local path where the downloaded file will be saved.
//	expectedSha256Sum: The expected SHA256 checksum of the downloaded file.
//	updateFunction: A function that takes a file path as input and returns an exit status code and an error.
//
// The updateDevice function performs the following steps:
//  1. Downloads the file from the remote server to the local path.
//  2. Verifies the downloaded file using its SHA256 checksum.
//  3. Calls the provided update function to update the device.
//  4. Logs the outcome of the update process, including any errors or status codes.
//
// It handles errors that may occur during the download, verification, or update process.
func updateDevice(remotePath, localPath, expectedSha256Sum string, updateFunction func(string) (int, error)) error {
	if updateFunction == nil {
		return errors.New("update function is nil")
	}

	// Download the file
	err := downloadFile(Handler.sftp, remotePath, localPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := verifyFile(localPath, expectedSha256Sum); err == nil {
		statusCode, err := updateFunction(localPath)
		if err != nil {
			log.Printf("Error updating device: %v. Status code: %d", err, statusCode)
		} else {
			if statusCode != 0 {
				log.Printf("Device update encountered an issue. Status code: %d", statusCode)
			} else {
				Logger.Infoln("Update successful!")
			}
		}
	} else {
		log.Printf("Update failed: %v", err)
	}
	return nil
}

// verifyFile verifies the integrity of a file by comparing its SHA-256 hash with an expected hash value.
// If the file's hash matches the expected hash, the function returns nil, indicating a successful verification.
// If there are any errors during the verification or if the file's hash does not match the expected hash,
// the function logs appropriate error messages and attempts to remove the file.
//
// Parameters:
//
//   - path:           The path to the file to be verified.
//   - expectedSha256Sum: The expected SHA-256 hash value to compare with the file's hash.
//
// Returns:
//
//	An error if any verification error occurs, or nil if the file's integrity is successfully verified.
func verifyFile(path, expectedSha256Sum string) error {
	fileSha256Sum, err := calculateSHA256Sum(path)
	if err != nil {
		log.Printf("Error calculating SHA-256 hash for file %s: %v", path, err)
		return err
	} else if expectedSha256Sum != fileSha256Sum {
		log.Printf("Error: The file %s has a different SHA-256 hash than expected. Removing it...", path)
		if removeErr := removeFile(path); removeErr != nil {
			log.Printf("Error removing file %s: %v", path, removeErr)
		} else {
			log.Printf("File %s has been removed successfully.", path)
		}
		return fmt.Errorf("corrupted file: %s", path)
	}
	Logger.Debugf("File %s has been successfully verified.", path)
	return nil
}

// downloadFile downloads a file from an SFTP server and saves it locally.
// It establishes an SSH connection using the provided SFTP configuration,
// downloads the specified remote file, and saves it to the local file path.
//
// Parameters:
//
//   - sftpHandler: An SFTP configuration containing host, port, username, and password.
//   - sremotePath: The remote file path on the SFTP server.
//   - localPath: The local file path where the downloaded file will be saved.
//
// Returns:
//
//   - An error if any connection or download error occurs; otherwise, it returns nil.
func downloadFile(sftpHandler *SFTPConfig, remotePath, localPath string) error {
	// Establish an SSH connection to the SFTP server
	sshClient, err := connectSSH(sftpHandler)
	if err != nil {
		Logger.Errorf("Error in SFTP %v\n", err)
		return err
	}
	defer sshClient.Close()

	// Download the file from the SFTP server and save it locally
	err = downloadFileSFTP(sshClient, remotePath, localPath)
	if err != nil {
		return err
	}

	return nil
}

// removeFile removes a file at the specified path.
func removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// NewSFTPConfig creates a new SFTPConfig with the provided parameters.
func newSFTPConfig(host, port, username, password string) *SFTPConfig {
	return &SFTPConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// ConnectSSH establishes an SSH connection to the SFTP server.
func connectSSH(config *SFTPConfig) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.host, config.port), sshConfig)
}

// downloadFileSFTP downloads a file from the SFTP server and saves it locally.
func downloadFileSFTP(sshClient *ssh.Client, remoteFilePath, localFilePath string) error {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client session: %v", err)
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	remoteFileInfo, err := remoteFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get remote file info: %v", err)
	}
	remoteFileSize := remoteFileInfo.Size()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	progressChan := make(chan float64)
	go func() {
		for progress := range progressChan {
			updateProgress(localFilePath, progress)
		}
	}()
	defer close(progressChan)

	copiedBytes, err := copyWithProgress(localFile, remoteFile, remoteFileSize, progressChan)
	if err != nil {
		return fmt.Errorf("download error: %v", err)
	}

	if copiedBytes != remoteFileSize {
		return fmt.Errorf("incomplete download: expected %d bytes, got %d bytes", remoteFileSize, copiedBytes)
	}

	Logger.Infof("Download of %s completed.", localFilePath)
	return nil
}

// calculateSHA256Sum calcula o hash SHA-256 de um arquivo e o retorna como uma string hexadecimal.
func calculateSHA256Sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	hashSum := hasher.Sum(nil)

	hashHex := hex.EncodeToString(hashSum)

	return hashHex, nil
}

// applyHlk7628Update updates a device using the sysupgrade command with the provided file path.
//
// Parameters:
//
//   - path: The path to the upgrade file.
//
// Returns:
//
//	int: The exit status code of the sysupgrade command.
//	error: An error, if any, during the execution of the command.
func applyHlk7628Update(path string) (int, error) {
	if !initializer.IsHlk7628UpdateEnabled() {
		return 0, nil
	}

	callHLK7628Events("updating")
	time.Sleep(3 * time.Second) // Wait to publish message
	cmd_string := "sysupgrade"
	args := []string{"-n", path} // Use -n to preserve configuration

	cmd := exec.Command(cmd_string, args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		callHLK7628Events("update fail")
		Logger.Errorf("Error executing the command: %v", err)
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}

// applyStm32Update updates a device using a custom script with the provided file path.
//
// Parameters:
//
//   - path: The path to the upgrade file.
//
// Returns:
//
//	int: The exit status code of the script execution.
//	error: An error, if any, during the execution of the script.
func applyStm32Update(path string) (int, error) {
	if !initializer.IsStm32UpdateEnabled() {
		return 0, nil
	}
	callSTM32Events("updating")
	cmd_string := "/opt/gabriel/bin/flash_stm32.sh"
	args := []string{path}

	write_success := false
	var err error = nil
	var output []byte
	var cmd *exec.Cmd
	for trials := 0; trials < 3 && !write_success; trials++ {
		Logger.Infof("Trying to write STM32 FW [%d]", trials)
		cmd = exec.Command(cmd_string, args...)
		output, err = cmd.CombinedOutput()
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "(100.00%) Done.") {
				write_success = true
				break
			}
		}
	}
	if !write_success {
		callSTM32Events("update fail")
		err = fmt.Errorf("Could not flash STM32 firmware through its entirety")
		Logger.Errorln("Failed")
	}
	callSTM32Events("updated")
	if err != nil {
		Logger.Errorf("Error executing the command: %v", err)
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}

// copyWithProgress copies data from a source reader to a destination writer,
// while tracking the copy progress.
//
// Parameters:
//
//	dst          - The destination writer where data will be copied to.
//	src          - The source reader from which data will be read.
//	size         - The total size of data to be copied (for progress calculation).
//	progressChan - A channel to send the copy progress (in percentage).
//
// Returns:
//
//	int64 - The total number of bytes copied.
//	error - Any error that occurs during the copy.
//
// The function reads data from the source reader and writes it to the destination
// writer in blocks, while updating the progress on the progressChan channel.
func copyWithProgress(dst io.Writer, src io.Reader, size int64, progressChan chan<- float64) (int64, error) {
	const bufferSize = 2097152

	buffer := make([]byte, bufferSize)
	var totalBytesCopied int64

	for {
		n, err := src.Read(buffer)
		if n > 0 {
			totalBytesCopied += int64(n)

			progress := float64(totalBytesCopied) / float64(size) * 100

			progressChan <- progress
		}

		if err != nil {
			if err == io.EOF {
				_, err = dst.Write(buffer[:n])
				if err != nil {
					return totalBytesCopied, err
				}
				break
			}
			return totalBytesCopied, err
		}

		_, err = dst.Write(buffer[:n])
		if err != nil {
			return totalBytesCopied, err
		}
	}

	return totalBytesCopied, nil
}

// updateProgress prints the download progress as a progress bar on the same line.
func updateProgress(localFilePath string, progress float64) {
	const progressBarWidth = 40 // Largura da barra de progresso em caracteres
	progressBar := strings.Repeat("=", int(progress*progressBarWidth/100))
	spaces := strings.Repeat(" ", progressBarWidth-int(progress*progressBarWidth/100))
	Logger.Infof("\rDownload of %s <%s%s> %.2f%%", localFilePath, progressBar, spaces, progress)
	if progress >= 100 {
		fmt.Println()
	}
}
