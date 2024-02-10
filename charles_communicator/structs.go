package charles_communicator

import (
	"container/list"

	"github.com/tarm/serial"
)

type CharlesMessage struct {
	version             uint8
	messageId           uint16
	messageType         uint8
	command             uint8
	dataLen             uint8
	data                string
	messageTimeout      int64
	messageCreationTime int64
	reqFunc             pointerToCharlesFunction
	externalData        interface{}
}

type CharlesCommunicatorHandler struct {
	respFunctionsList       list.List
	messagesToSendList      list.List
	messagesWaitingResponse list.List
	rcvMsgState             int
	config                  *serial.Config
	port                    *serial.Port
	ValidPort               bool
	rcvMsgBufferControl     uint
	rcvMsgBuffer            [BUFFER_SIZE]byte
	msgIdControl            uint16
}

type RespFuction struct {
	respFunction pointerToCharlesFunction
	externalData interface{}
	messageType  uint8
	command      uint8
}

type respStatusParameters struct {
	messageType uint8
	message     string
}
