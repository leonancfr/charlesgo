package utils

import (
	"common"
	"crypto/aes"
	"device_info"
	"gablogger"
	"peripherals"
)

var Logger = gablogger.Logger()

func Banner() {
	deviceId, err := device_info.GetDeviceId()

	if err != nil {
		deviceId = "Undefined"
	}

	Logger.Debug("=============")
	Logger.Debug("Version ", common.VERSION)
	Logger.Debug("Device mac address: ", deviceId)
	os_version, _ := device_info.GetOSVersion()
	Logger.Debug("Device OSVersion: ", os_version)
	fw_version, err := peripherals.GetFirmwareVersion()
	if err != nil {
		fw_version = err.Error()
	}
	Logger.Debug("Device FWVersion: ", fw_version)
	Logger.Debug("=============")
}

// EncryptAES256ECB takes a key and data, and returns the encrypted data using AES-256-ECB.
func EncryptAES256ECB(key []byte, data []byte) ([]byte, error) {
	// Create a new AES cipher block with the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Add padding to the data if necessary
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	if padding > 0 {
		for i := 0; i < padding; i++ {
			data = append(data, byte(padding))
		}
	}

	// Initialize a byte slice to store the encrypted data
	ciphertext := make([]byte, len(data))
	// Encrypt each 16-byte block using ECB
	for i := 0; i < len(data); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], data[i:i+aes.BlockSize])
	}

	return ciphertext, nil
}

func GetUserName() string {
	deviceId, err := device_info.GetDeviceId()
	if err != nil {
		return ""
	}
	return "router_" + deviceId
}
