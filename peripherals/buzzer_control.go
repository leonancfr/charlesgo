package peripherals

import (
	"charles_communicator"
	"event_control"
)

func setupBuzzerControl() {
	event_control.RegisterToReceiveEvent(GetTamperEventId(), activateBuzzerByTamper, nil)
}

func activateBuzzerByTamper(messageType, command uint8, message string, externalData interface{}) {
	if message == "Open" {
		resp, err := CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_BUZZER_ENABLE, "10", charles_communicator.WAIT_MESSAGE_RESPONSE_TIMEOUT)
		if err != nil {
			Logger.Errorln("Cannot send message to enable buzzer.", err)
			event_control.CallRegisteredEventFunctions(GetBuzzerEventId(), charles_communicator.MSG_TYPE_ERROR, charles_communicator.MSG_CMD_BUZZER_ENABLE, err.Error())
		} else {
			event_control.CallRegisteredEventFunctions(GetBuzzerEventId(), charles_communicator.MSG_TYPE_RESP, charles_communicator.MSG_CMD_BUZZER_ENABLE, resp)
		}
	}

	if message == "Close" {
		resp, err := CCHandler.SendMessage(charles_communicator.MSG_TYPE_SET, charles_communicator.MSG_CMD_BUZZER_DISABLE, "", 3000)
		if err != nil {
			Logger.Errorln("Cannot send message to disable buzzer.", err)
			event_control.CallRegisteredEventFunctions(GetBuzzerEventId(), charles_communicator.MSG_TYPE_ERROR, charles_communicator.MSG_CMD_BUZZER_DISABLE, err.Error())
		} else {
			event_control.CallRegisteredEventFunctions(GetBuzzerEventId(), charles_communicator.MSG_TYPE_RESP, charles_communicator.MSG_CMD_BUZZER_DISABLE, resp)
		}
	}
}
