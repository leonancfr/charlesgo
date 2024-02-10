package peripherals

import (
	"charles_communicator"
	"gablogger"
)

var Logger = gablogger.Logger()

var CCHandler *charles_communicator.CharlesCommunicatorHandler

func Setup(ccHandler *charles_communicator.CharlesCommunicatorHandler) {
	CCHandler = ccHandler

	CCHandler.RegisterFunctionToRcvMsg(charles_communicator.MSG_TYPE_GET, charles_communicator.MSG_CMD_GET_WATCHDOG, respGetWatchdog, nil)
	CCHandler.RegisterFunctionToRcvMsg(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_TAMPER_EVENT, respSetTamperEvent, nil)
	CCHandler.RegisterFunctionToRcvMsg(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_BUZZER_DISABLE, respBuzzerEvent, nil)
	CCHandler.RegisterFunctionToRcvMsg(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_POWER_SOURCE, respPowerSourceEvent, nil)

	setupBuzzerControl()
}
