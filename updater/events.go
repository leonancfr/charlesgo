package updater

import (
	"event_control"
)

var hlk7628UpdateEventId int
var stm32UpdateEventId int

func GetHLK7628UpdateEventId() int {
	if hlk7628UpdateEventId == 0 {
		hlk7628UpdateEventId = event_control.CreateEventId()
	}
	return hlk7628UpdateEventId
}

func GetSTM32pdateEventId() int {
	if stm32UpdateEventId == 0 {
		stm32UpdateEventId = event_control.CreateEventId()
	}
	return stm32UpdateEventId
}

func callSTM32Events(message string) {
	event_control.CallRegisteredEventFunctions(GetSTM32pdateEventId(), 0, 0, message)
}

func callHLK7628Events(message string) {
	event_control.CallRegisteredEventFunctions(GetHLK7628UpdateEventId(), 0, 0, message)
}
