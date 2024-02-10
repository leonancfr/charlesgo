package monitor

import "charles_communicator"

func sendTamperEvent(messageType, command uint8, message string, externalData interface{}) {
	publishMetric(topicTamperStatus, message)
}

func sendPowerSourceEvent(messageType, command uint8, message string, externalData interface{}) {
	publishMetric(topicPowerSource, message)
}

func sendBuzzerEvent(messageType, command uint8, message string, externalData interface{}) {
	if command == charles_communicator.MSG_CMD_BUZZER_ENABLE {
		switch messageType {
		case charles_communicator.MSG_TYPE_RESP:
			publishMetric(topicBuzzerStatus, "active")
		default:
			publishMetric(topicBuzzerStatus, "activation error")
		}
	}

	if command == charles_communicator.MSG_CMD_BUZZER_DISABLE {
		switch messageType {
		case charles_communicator.MSG_TYPE_RESP:
			publishMetric(topicBuzzerStatus, "inactive")
		case charles_communicator.MSG_TYPE_SET:
			publishMetric(topicBuzzerStatus, "inactive")
		default:
			publishMetric(topicBuzzerStatus, "deactivation error")
		}
	}
}

func sendUpdateSTM32Event(messageType, command uint8, message string, externalData interface{}) {
	publishMetric(topicUpdateSTM32, message)
}

func sendUpdateHLK7628Event(messageType, command uint8, message string, externalData interface{}) {
	publishMetric(topicUpdateHLK7628, message)
}
