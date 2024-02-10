package event_control

import "sync"

type pointerToEventFunction func(messageType, command uint8, message string, externalData interface{})

type eventDataStructure struct {
	function     pointerToEventFunction
	externalData interface{}
}

var (
	eventIdCounter         = 0
	eventIdCounterMutex    sync.Mutex
	eventFunctionsRegister = make(map[int][]eventDataStructure)
)

// RegisterToReceiveEvent registers a callback function to handle events.
// It associates the provided callback function and external data with the specified event.
func RegisterToReceiveEvent(event int, callbackFunction pointerToEventFunction, externalData interface{}) {
	eventFunctionsRegister[event] = append(eventFunctionsRegister[event], eventDataStructure{function: callbackFunction, externalData: externalData})
}

func CallRegisteredEventFunctions(event int, messageType, command uint8, message string) {
	if eventFunctions, ok := eventFunctionsRegister[event]; ok {
		for _, eventFunctionStruct := range eventFunctions {
			go eventFunctionStruct.function(messageType, command, message, eventFunctionStruct.externalData)
		}
	}
}

func CreateEventId() int {
	eventIdCounterMutex.Lock()
	defer eventIdCounterMutex.Unlock()
	eventIdCounter = eventIdCounter + 1

	return eventIdCounter
}
