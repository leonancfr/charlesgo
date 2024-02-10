package peripherals

import "charles_communicator"

func GetModemSignalStrength() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_MODEM_SIGNAL, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetModemConnectionType() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_MODEM_CONN_TYPE, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetModemConnectionBand() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_MODEM_CONN_BAND, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetSIMCardType() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_SIM_TYPE, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetSIMCardICCID() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_SIM_ICCID, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetSIMCardCarrier() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_SIM_CARRIER, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}
