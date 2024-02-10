package peripherals

import (
	"charles_communicator"
	"event_control"
)

var eventTamperId int
var eventWatchdogId int
var eventBuzzerId int
var eventPowerSourceId int

func respGetWatchdog(messageType, command uint8, message string, messageToken, externalData interface{}) {
	Logger.Debugln("Watchdog response message sent")
	CCHandler.SendRespMessage(messageToken.(*charles_communicator.CharlesMessage), "OK")
	event_control.CallRegisteredEventFunctions(GetWatchdogEventId(), messageType, command, message)
}

func respSetTamperEvent(messageType, command uint8, message string, messageToken, externalData interface{}) {
	CCHandler.SendRespMessage(messageToken.(*charles_communicator.CharlesMessage), "OK")
	event_control.CallRegisteredEventFunctions(GetTamperEventId(), messageType, command, message)
}

func respBuzzerEvent(messageType, command uint8, message string, messageToken, externalData interface{}) {
	CCHandler.SendRespMessage(messageToken.(*charles_communicator.CharlesMessage), "OK")
	event_control.CallRegisteredEventFunctions(GetBuzzerEventId(), messageType, command, message)
}

func respPowerSourceEvent(messageType, command uint8, message string, messageToken, externalData interface{}) {
	CCHandler.SendRespMessage(messageToken.(*charles_communicator.CharlesMessage), "OK")
	event_control.CallRegisteredEventFunctions(GetPowerSourceEventId(), messageType, command, message)
}

func GetTamperEventId() int {
	if eventTamperId == 0 {
		eventTamperId = event_control.CreateEventId()
	}
	return eventTamperId
}

func GetWatchdogEventId() int {
	if eventWatchdogId == 0 {
		eventWatchdogId = event_control.CreateEventId()
	}
	return eventWatchdogId
}

func GetBuzzerEventId() int {
	if eventBuzzerId == 0 {
		eventBuzzerId = event_control.CreateEventId()
	}
	return eventBuzzerId
}

func GetPowerSourceEventId() int {
	if eventPowerSourceId == 0 {
		eventPowerSourceId = event_control.CreateEventId()
	}
	return eventPowerSourceId
}
