package peripherals

import "charles_communicator"

func GetPowerSource() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_POWER_SOURCE, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetHasBMS() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_HAS_BMS, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetBatteryLevel() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_BATTERY_LEVEL, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}
