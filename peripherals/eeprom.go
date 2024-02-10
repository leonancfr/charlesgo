package peripherals

import "charles_communicator"

func setEepromDisableWriteProtection() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_EEPROM_DISABLE_WRITE_PROTECTION, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func setEepromEnableWriteProtection() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_EEPROM_ENABLE_WRITE_PROTECTION, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}
