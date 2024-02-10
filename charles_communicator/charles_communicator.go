package charles_communicator

import (
	"errors"
	"gablogger"
	"initializer"
	"io"
	"reflect"
	"time"

	"github.com/tarm/serial"
)

var Logger = gablogger.Logger()

type pointerToCharlesFunction func(messageType, command uint8, message string, messageToken, externalData interface{})

var CCHandler *CharlesCommunicatorHandler

func GetCharlesCommunicatorHandler() *CharlesCommunicatorHandler {
	return CCHandler
}

func InitCC() {
	CCHandler = &CharlesCommunicatorHandler{
		config:      &serial.Config{Name: "/dev/ttyS1", Baud: 115200, ReadTimeout: time.Millisecond * 500},
		rcvMsgState: WAITING_FOR_MESSAGE,
	}
	OpenPort()
}

func ClosePort() {
	Logger.Debugln("Closing port", CCHandler.config.Name)

	CCHandler.ValidPort = false
	CCHandler.port.Flush()
	err := CCHandler.port.Close()
	if err != nil {
		Logger.Infoln("Error closing serial port:", err)
	}
}

func OpenPort() bool {
	Logger.Debugln("Opening port", CCHandler.config.Name)

	port, err := serial.OpenPort(CCHandler.config)
	if err != nil {
		Logger.Errorf("Error opening port %s: %s", CCHandler.config.Name, err)
		return false
	}
	CCHandler.port = port
	CCHandler.ValidPort = true
	return true
}

func (h *CharlesCommunicatorHandler) Start() {
	if !initializer.IsSupervisorEnable() {
		Logger.Debugln("Supervisor is disabled")
		return
	}
	Logger.Infoln("Charles Communicator Handler started!")
	var err error
	var n int
	buf := make([]byte, 1)
	for {
		for {
			if h.ValidPort {
				n, err = h.port.Read(buf)
				break
			} else {
				time.Sleep(500 * time.Millisecond)
			}
		}
		if err != nil && err != io.EOF {
			Logger.Errorf("Error in stm Handler: %v\n", err)
			ClosePort()
			time.Sleep(2 * time.Second)
			OpenPort()
		} else {
			if n > 0 {
				switch buf[0] {
				case '[':
					h.rcvMsgState = RECEIVING_DATA
					h.rcvMsgBufferControl = 0
				case ']':
					if h.rcvMsgState == RECEIVING_DATA {
						h.rcvMsgBuffer[h.rcvMsgBufferControl] = byte('\x00')
						h.rcvMsgBufferControl++
						h.processMessage()
					}
					h.rcvMsgState = WAITING_FOR_MESSAGE
				default:
					if h.rcvMsgState == RECEIVING_DATA {
						if h.rcvMsgBufferControl >= (BUFFER_SIZE - 1) {
							h.rcvMsgState = WAITING_FOR_MESSAGE
						} else {
							h.rcvMsgBuffer[h.rcvMsgBufferControl] = buf[0]
							h.rcvMsgBufferControl++
						}
					}
				}
			} else {
				h.checkWaitingTimeouts()
				h.sendMessageIfExists()
			}
		}
	}
}

func (h *CharlesCommunicatorHandler) processMessage() {
	rawMessage := string(h.rcvMsgBuffer[:])
	rawMessage = rawMessage[:(h.rcvMsgBufferControl - 1)]
	message := decodeCharlesMessage(rawMessage)
	if message != nil {
		switch message.messageType {
		case MSG_TYPE_GET:
			fallthrough
		case MSG_TYPE_SET:
			if !h.callRespFunc(message) {
				h.SendErrorMessage(message, "unsupported command")
			}
		case MSG_TYPE_RESP:
			fallthrough
		case MSG_TYPE_ERROR:
			h.callWaitingRespFunc(message)
		}
	}
}

func (h *CharlesCommunicatorHandler) callRespFunc(message *CharlesMessage) bool {
	for e := h.respFunctionsList.Front(); e != nil; e = e.Next() {
		respFunction, ok := e.Value.(*RespFuction)
		if ok {
			if respFunction.command == message.command && respFunction.messageType == message.messageType {
				if respFunction.respFunction != nil {
					respFunction.respFunction(message.messageType, message.command, message.data, message, respFunction.externalData)
				}
				return true
			}
		}
	}
	return false
}

func (h *CharlesCommunicatorHandler) callWaitingRespFunc(message *CharlesMessage) bool {

	for e := h.messagesWaitingResponse.Front(); e != nil; e = e.Next() {
		messageWaiting, ok := e.Value.(*CharlesMessage)
		if ok {

			if message.messageId == messageWaiting.messageId {

				if messageWaiting.reqFunc != nil {
					messageWaiting.reqFunc(message.messageType, message.command, message.data, message, messageWaiting.externalData)
				}
				h.messagesWaitingResponse.Remove(e)
				return true
			}
		}
	}
	return false
}

func (h *CharlesCommunicatorHandler) SendErrorMessage(messageToken *CharlesMessage, data string) {
	message := CharlesMessage{
		version:             messageToken.version,
		messageId:           messageToken.messageId,
		messageType:         MSG_TYPE_ERROR,
		command:             messageToken.command,
		dataLen:             uint8(len(data) + 1),
		data:                data,
		messageTimeout:      5000,
		messageCreationTime: time.Now().UnixMilli(),
	}
	h.registerMessage(&message)
}

func (h *CharlesCommunicatorHandler) SendRespMessage(messageToken *CharlesMessage, data string) {
	message := CharlesMessage{
		version:             messageToken.version,
		messageId:           messageToken.messageId,
		messageType:         MSG_TYPE_RESP,
		command:             messageToken.command,
		dataLen:             uint8(len(data) + 1),
		data:                data,
		messageTimeout:      5000,
		messageCreationTime: time.Now().UnixMilli(),
	}
	h.registerMessage(&message)
}

func (h *CharlesCommunicatorHandler) SendGetMessage(command uint8, data string, timeout int64, toRcvFunction pointerToCharlesFunction, externalData interface{}) error {
	messageId, err := h.requestMessageId()
	if err != nil {
		return errors.New("cannot send message")
	}

	message := CharlesMessage{
		version:             PROTOCOL_VERSION,
		messageId:           messageId,
		messageType:         MSG_TYPE_GET,
		command:             command,
		dataLen:             uint8(len(data) + 1),
		data:                data,
		messageTimeout:      timeout,
		messageCreationTime: time.Now().UnixMilli(),
		reqFunc:             toRcvFunction,
		externalData:        externalData,
	}

	h.registerMessage(&message)

	return nil
}

func (h *CharlesCommunicatorHandler) SendSetMessage(command uint8, data string, timeout int64, toRcvFunction pointerToCharlesFunction, externalData interface{}) error {
	messageId, err := h.requestMessageId()
	if err != nil {
		return errors.New("cannot send message")
	}

	message := CharlesMessage{
		version:             PROTOCOL_VERSION,
		messageId:           messageId,
		messageType:         MSG_TYPE_SET,
		command:             command,
		dataLen:             uint8(len(data) + 1),
		data:                data,
		messageTimeout:      timeout,
		messageCreationTime: time.Now().UnixMilli(),
		reqFunc:             toRcvFunction,
		externalData:        externalData,
	}

	h.registerMessage(&message)

	return nil
}

func (h *CharlesCommunicatorHandler) RegisterFunctionToRcvMsg(messageType, command uint8, respFunction pointerToCharlesFunction, externalData interface{}) {
	newFunction := RespFuction{
		respFunction: respFunction,
		externalData: externalData,
		messageType:  messageType,
		command:      command,
	}
	h.respFunctionsList.PushBack(&newFunction)
}

func (h *CharlesCommunicatorHandler) registerMessage(message *CharlesMessage) {
	h.messagesToSendList.PushBack(message)
}

func (h *CharlesCommunicatorHandler) requestMessageId() (uint16, error) {
	lastMessageId := h.msgIdControl
	usedMessageIds := h.listAllUsedMessageIds()
	for {
		pretendedMessageId := 2*(h.msgIdControl) + 1
		h.msgIdControl++

		if !usedMessageIds[pretendedMessageId] {
			return pretendedMessageId, nil
		}

		if lastMessageId == pretendedMessageId {
			return 0, errors.New("no message id available")
		}
	}
}

func (h *CharlesCommunicatorHandler) listAllUsedMessageIds() map[uint16]bool {
	listOfUsedMessageIds := make(map[uint16]bool)

	for e := h.messagesToSendList.Front(); e != nil; e = e.Next() {
		message, ok := e.Value.(CharlesMessage)
		if ok {
			listOfUsedMessageIds[message.messageId] = true
		}
	}

	for e := h.messagesWaitingResponse.Front(); e != nil; e = e.Next() {
		message, ok := e.Value.(CharlesMessage)
		if ok {
			listOfUsedMessageIds[message.messageId] = true
		}
	}

	return listOfUsedMessageIds
}

func (h *CharlesCommunicatorHandler) sendMessageIfExists() {
	message, err := h.getMessageToSend()
	if err == nil {
		if message != nil {
			message.messageCreationTime = time.Now().UnixMilli()

			serializedMessage, err := encodeCharlesMessage(message)
			if err == nil {
				for {
					if h.ValidPort {
						_, err = h.port.Write([]byte(serializedMessage))
						break
					} else {
						time.Sleep(500 * time.Millisecond)
					}
				}
				if err == nil {
					switch message.messageType {
					case MSG_TYPE_GET:
						fallthrough
					case MSG_TYPE_SET:
						h.messagesWaitingResponse.PushBack(message)
					}
				} else {
					Logger.Errorln("Error writing to port", err)
				}
			}
			if err != nil {
				Logger.Errorln("Error sending message:", err)
			}
		}
	} else {
		Logger.Errorln("Error requesting a message to send:", err)
	}
}

func (h *CharlesCommunicatorHandler) getMessageToSend() (*CharlesMessage, error) {
	if h.messagesToSendList.Len() > 0 {
		listElement := h.messagesToSendList.Back()
		message, ok := listElement.Value.(*CharlesMessage)
		h.messagesToSendList.Remove(listElement)
		if ok {
			return message, nil
		} else {
			return nil, errors.New("cannot decode element of message list, removing it")
		}
	}
	return nil, nil
}

func (h *CharlesCommunicatorHandler) checkWaitingTimeouts() {
	actualTime := time.Now().UnixMilli()

	for e := h.messagesWaitingResponse.Front(); e != nil; e = e.Next() {
		message, ok := e.Value.(*CharlesMessage)
		if ok {
			if (actualTime - message.messageCreationTime) > message.messageTimeout {
				if message.reqFunc != nil {
					message.reqFunc(MSG_TYPE_TIMEOUT, message.command, "", message, message.externalData)
				}
				h.messagesWaitingResponse.Remove(e)
			}
		}
	}
}

func (h *CharlesCommunicatorHandler) SendMessage(messageType, messageCommand uint8, messageString string, messageTimeout int64) (string, error) {
	if !initializer.IsSupervisorEnable() {
		return "", errors.New("supervisor is disabled")
	}
	responseChannel := make(chan respStatusParameters)

	switch messageType {
	case MSG_TYPE_GET:
		h.SendGetMessage(messageCommand, "", messageTimeout, respSetGet, responseChannel)
	case MSG_TYPE_SET:
		h.SendSetMessage(messageCommand, messageString, messageTimeout, respSetGet, responseChannel)
	default:
		return "", errors.New("message type not available")
	}
	response := <-responseChannel

	return returnStringDataAndError(response)
}

func returnStringDataAndError(responseStruct respStatusParameters) (string, error) {
	responseData := responseStruct.message
	responseType := responseStruct.messageType

	if responseType == MSG_TYPE_ERROR {
		return "", errors.New(responseData)
	}
	if responseData == "" {
		return "", errors.New("timeout")
	}
	return responseData, nil
}

func respSetGet(messageType, command uint8, message string, messageToken, externalData interface{}) {
	ch, cond := externalData.(chan respStatusParameters)

	if !cond {
		Logger.Errorln("expected respStatusParamaters channel, received ", reflect.TypeOf(externalData))
		return
	}

	ch <- respStatusParameters{messageType, message}
}
