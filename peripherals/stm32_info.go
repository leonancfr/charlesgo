package peripherals

import (
	"charles_communicator"
	"errors"
)

func GetSTM32Temperature() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_STM32_TEMPERATURE, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetFirmwareVersion() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_FIRMWARE_VERSION, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func SetOsVersion(osVersion string) (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_OS_VERSION, osVersion, charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func SetSerialNumber(serialNumber string) (string, error) {
	if _, err := setEepromDisableWriteProtection(); err != nil {
		return "", errors.New("unable to disable write protection")
	}

	data, erro := CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_SERIAL_NUMBER, serialNumber, charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)

	if _, err := setEepromEnableWriteProtection(); err != nil {
		return "", errors.New("unable to enable write protection")
	}

	return data, erro
}
