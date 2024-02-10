package socketxp

import (
	"encoding/json"
	"errors"
	"fmt"
	"gablogger"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var Logger = gablogger.Logger()
var Handler *SocketXp

type SocketXp struct {
	configJSONPath  string
	deviceKeyPath   string
	CredentialTopic string
	workDir         string
}

type tunnel struct {
	Destination  string `json:"destination"`
	CustomDomain string `json:"custom_domain"`
	Subdomain    string `json:"subdomain"`
}

type configJson struct {
	WorkDir string   `json:"work_dir"`
	Tunnels []tunnel `json:"tunnels"`
}

type deviceKey struct {
	DeviceId    string `json:"DeviceId"`
	DeviceKey   string `json:"DeviceKey"`
	DeviceName  string `json:"DeviceName"`
	DeviceGroup string `json:"DeviceGroup"`
}

func InitSocketXP() {
	workDir := "/opt/socketxp"
	username := utils.GetUserName()
	Handler = &SocketXp{
		workDir:         workDir,
		configJSONPath:  fmt.Sprintf("%s/config.json", workDir),
		deviceKeyPath:   fmt.Sprintf("%s/device.key", workDir),
		CredentialTopic: fmt.Sprintf("devices/%s/provisioning/socketxp", username),
	}
}

func UpdateCredentialsCallback(client mqtt.Client, message mqtt.Message) {
	receivedKey, err := decodeMessage(string(message.Payload()))
	if err != nil {
		Logger.Errorf("Cannot decode message. %v", err)
		return
	}

	newConfigJson := Handler.createConfigJson(receivedKey.DeviceName)

	actualConfigJson, err := Handler.loadConfigJsonFile()
	if err != nil {
		Logger.Warningf("Cannot load config json file. %v", err)
	}

	actualDeviceKey, err := Handler.loadDeviceKeyFile()
	if err != nil {
		Logger.Warningf("Cannot load device key file. %v", err)
	}

	if !actualConfigJson.isEqual(newConfigJson) || !actualDeviceKey.isEqual(receivedKey) {
		if Handler.updateCredentialsFiles(receivedKey, newConfigJson) != nil {
			return
		}
		Logger.Infoln("SocketXP credentials files have been updated.")
		if restartSocketXp() != nil {
			return
		}
		Logger.Infoln("SocketXP have been updated.")
	}
}

func (s SocketXp) createConfigJson(DeviceName string) configJson {
	return configJson{
		WorkDir: s.workDir,
		Tunnels: []tunnel{
			{
				Destination: "tcp://127.0.0.1:2202",
			},
			{
				Destination:  "tcp://127.0.0.1:4404",
				CustomDomain: "", Subdomain: "gabriel-tech-" + DeviceName,
			},
		},
	}
}

func (s SocketXp) updateCredentialsFiles(newDeviceKey deviceKey, newConfigJson configJson) error {

	if err := s.createWorkDir(); err != nil {
		Logger.Errorf("Cannot create work dir %s. %v.", s.workDir, err)
		return err
	}

	if err := s.createConfigJsonFile(newConfigJson); err != nil {
		Logger.Errorf("Cannot update config json file %v", err)
		return err
	}

	if err := s.createDeviceKeyFile(newDeviceKey); err != nil {
		Logger.Errorf("Cannot update device key file %v", err)
		return err
	}

	return nil
}

func restartSocketXp() error {
	cmd := exec.Command("service", "socket_xp", "restart")
	err := cmd.Run()
	if err != nil {
		Logger.Errorf("Error: Failed to restart the SocketXP service. %v", err)
	}
	return err
}

func decodeMessage(rawMessage string) (deviceKey, error) {
	var key deviceKey
	err := json.Unmarshal([]byte(rawMessage), &key)
	if err != nil {
		Logger.Errorln("Cannot parse message:", err)
	} else if key.DeviceId == "" || key.DeviceKey == "" || key.DeviceName == "" {
		err = errors.New("insufficient data received")
	}

	return key, err
}

func (s SocketXp) loadDeviceKeyFile() (deviceKey, error) {
	data, err := os.ReadFile(s.deviceKeyPath)
	if err != nil {
		return deviceKey{}, err
	}

	var key deviceKey
	if err := json.Unmarshal(data, &key); err != nil {
		return deviceKey{}, err
	}

	return key, nil
}

func (s SocketXp) loadConfigJsonFile() (configJson, error) {

	data, err := os.ReadFile(s.configJSONPath)
	if err != nil {
		return configJson{}, err
	}

	var config configJson
	err = json.Unmarshal(data, &config)
	if err != nil {
		return configJson{}, err
	}

	return config, nil
}

func (s deviceKey) isEqual(value deviceKey) bool {
	return s.DeviceId == value.DeviceId &&
		s.DeviceKey == value.DeviceKey &&
		s.DeviceName == value.DeviceName &&
		s.DeviceGroup == value.DeviceGroup
}

func (c configJson) isEqual(value configJson) bool {
	if c.WorkDir != value.WorkDir {
		return false
	}

	if len(c.Tunnels) != len(value.Tunnels) {
		return false
	}

	for i := range c.Tunnels {
		if c.Tunnels[i] != value.Tunnels[i] {
			return false
		}
	}

	return true
}

func (s SocketXp) createConfigJsonFile(newConfigJson configJson) error {
	jsonData, err := json.MarshalIndent(newConfigJson, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.configJSONPath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s SocketXp) createDeviceKeyFile(newDeviceKey deviceKey) error {
	deviceKeyJSON, err := json.MarshalIndent(newDeviceKey, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.deviceKeyPath, deviceKeyJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s SocketXp) createWorkDir() error {
	err := os.MkdirAll(s.workDir, os.ModePerm)
	return err
}

// IsConnected checks whether the connection to SocketXP is established and functioning correctly.
// It returns true if the connection is active and all checks are successful.
// Otherwise, it returns false along with an error indicating the reason for the failure.
func IsConnected() (string, error) {
	if Handler != nil {
		if !checkSocketXpFiles() {
			return "not provisioned", fmt.Errorf("SocketXp files not found")
		}

		if !checkSocketXpService() {
			return "stopped", fmt.Errorf("SocketXp service not running")
		}

		if !checkSocketXpConnectivity() {
			return "not working", fmt.Errorf("SocketXP connectivity not working")
		}
		return "connected", nil
	} else {
		return "monitoring error", fmt.Errorf("SocketXp handler not initialized")
	}
}

func checkSocketXpFiles() bool {
	return configJsonFileIsValid() && deviceKeyFileIsValid()
}

func checkSocketXpService() bool {
	cmd := exec.Command("service", "socket_xp", "status")

	output, err := cmd.CombinedOutput()
	if err != nil {
		Logger.WithField("socketxp", "service").Error("Failed to check SocketXP service status")
		return false
	}

	status := strings.TrimSpace(string(output))
	if status != "running" {
		Logger.WithField("socketxp", "service").Errorf("SocketXP service is not running (Status: %s)", status)
		return false
	}

	Logger.WithField("socketxp", "service").Debug("SocketXP service is running")
	return true
}

func checkSocketXpConnectivity() bool {
	deviceKey, err := Handler.loadDeviceKeyFile()
	if err != nil {
		Logger.WithField("socketxp", "connectivity").Error("Failed to load device key file")
		return false
	}

	autoPingURL := "https://gabriel-tech-" + deviceKey.DeviceName + ".socketxp.com"

	client := http.Client{Timeout: 5 * time.Second}

	response, err := client.Head(autoPingURL)
	if err != nil {
		Logger.WithField("socketxp", "connectivity").Errorf("Failed to perform request: %v", err)
		return false
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		Logger.WithField("socketxp", "connectivity").Debug("Connectivity check passed")
		return true
	} else {
		Logger.WithFields(logrus.Fields{
			"socketxp":    "connectivity",
			"status_code": response.StatusCode,
			"status":      response.Status,
		}).Error("Connectivity check failed")
		return false
	}
}

func configJsonFileIsValid() bool {
	configJson, err := Handler.loadConfigJsonFile()
	if err != nil {
		Logger.WithField("socketxp", "invalid-file").Error("Failed to parse config.json file")
		return false
	}

	if configJson.WorkDir != Handler.workDir {
		Logger.WithField("socketxp", "invalid-file").Error("Invalid work directory")
		return false
	}

	for _, tunnel := range configJson.Tunnels {
		if tunnel.Destination != "" {
			return true
		}
	}

	Logger.WithField("socketxp", "invalid-file").Error("No tunnels registered")
	return false
}

func deviceKeyFileIsValid() bool {
	deviceKey, err := Handler.loadDeviceKeyFile()
	if err != nil {
		Logger.WithField("socketxp", "invalid-file").Error("Failed to parse device.key file")
		return false
	}

	if deviceKey.DeviceName == "" || deviceKey.DeviceId == "" || deviceKey.DeviceKey == "" {
		Logger.WithField("socketxp", "invalid-file").Error("Incomplete or missing fields in device.key")
		return false
	}

	return true
}
