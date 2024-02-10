package device_info

import (
	"bufio"
	"common"
	"fmt"
	"initializer"
	"io"
	"os"
	"regexp"
	"strings"
)

// getValueFromFile reads a specified key's value from a given file.
func getValueFromFile(path, key string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regex := regexp.MustCompile(key + `="?([^"]+)"?`)

	for scanner.Scan() {
		line := scanner.Text()

		if match := regex.FindStringSubmatch(line); len(match) > 1 {
			return match[1], nil
		}
	}
	return "", fmt.Errorf("%s not found in the file", key)
}

func GetOSVersion() (string, error) {
	return getValueFromFile("/etc/os-release", "VERSION")
}

func GetCharlesGoVersion() (string, error) {
	return common.VERSION, nil
}

func findMtdIndex(name string) (string, error) {
	data, err := os.ReadFile("/proc/mtd")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "\""+name+"\"") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				part := parts[0]
				index := strings.TrimPrefix(part, "mtd")
				return index, nil
			}
		}
	}

	return "", fmt.Errorf("MTD partition not found: %s", name)
}

func findMtdPart(name string) (string, error) {
	index, err := findMtdIndex(name)
	if err != nil {
		return "", err
	}

	prefix := "/dev/mtdblock"
	if _, err := os.ReadDir("/dev/mtdblock"); err == nil {
		prefix = "/dev/mtdblock/"
	}

	if index != "" {
		return prefix + index, nil
	}

	return "", fmt.Errorf("MTD partition not found: %s", name)
}

func getMACBinary(path string, offset int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("get_mac_binary: file %s not found", path)
	}
	defer file.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	macBytes := make([]byte, 6)
	_, err = io.ReadFull(file, macBytes)
	if err != nil {
		return "", err
	}

	macString := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		macBytes[0], macBytes[1], macBytes[2],
		macBytes[3], macBytes[4], macBytes[5])

	return macString, nil
}

func GetDeviceId() (string, error) {
	if initializer.GetLabel() != "" {
		return initializer.GetLabel(), nil
	}
	return getMac("factory", int64(0x4))
}

func GetSecondaryMAC() (string, error) {
	return getMac("factory", int64(0x2e))
}

func getMac(mtdName string, offset int64) (string, error) {
	partition, err := findMtdPart(mtdName)
	if err != nil {
		return "", err
	}

	mac, err := getMACBinary(partition, offset)
	if err != nil {
		return "", err
	}

	return mac, nil
}
