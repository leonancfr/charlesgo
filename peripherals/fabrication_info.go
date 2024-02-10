package peripherals

import "charles_communicator"

func GetPCBReview() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_PCB_REV, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}

func GetPCBBatch() (string, error) {
	return CCHandler.SendMessage(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_BATCH_NUMBER, "", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
}
