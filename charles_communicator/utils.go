package charles_communicator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func decodeCharlesMessage(message string) *CharlesMessage {
	decodedMessage := &CharlesMessage{}

	fields := map[string]interface{}{
		"version":    &decodedMessage.version,
		"type":       &decodedMessage.messageType,
		"command":    &decodedMessage.command,
		"message_id": &decodedMessage.messageId,
		"data_len":   &decodedMessage.dataLen,
	}

	parts := strings.Split(message, ";")
	for _, part := range parts {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			return nil
		}
		key := kv[0]
		value := kv[1]
		if ptr, ok := fields[key]; ok {
			if intValue, err := strconv.Atoi(value); err == nil {
				switch v := ptr.(type) {
				case *uint8:
					*v = uint8(intValue)
				case *uint16:
					*v = uint16(intValue)
				}
			} else {
				return nil
			}
		} else if key == "data" {
			decodedMessage.data = value
		}
	}
	return decodedMessage
}

func encodeCharlesMessage(message *CharlesMessage) (string, error) {
	serializedMessage := fmt.Sprintf("[version:%d;type:%d;command:%d;message_id:%d;data_len:%d;data:%s]",
		message.version, message.messageType, message.command, message.messageId, message.dataLen, message.data)

	if len(serializedMessage)+1 > BUFFER_SIZE {
		return "", errors.New("message is too long")
	}

	return serializedMessage, nil
}

func TypeToString(value uint8) string {
	switch value {
	case MSG_TYPE_ERROR:
		return "ERROR"
	case MSG_TYPE_RESP:
		return "RESPONSE"
	case MSG_TYPE_TIMEOUT:
		return "TIMEOUT"
	case MSG_TYPE_GET:
		return "GET"
	case MSG_TYPE_SET:
		return "SET"
	default:
		return "undefined"
	}
}

func CommandToString(value uint8) string {
	switch value {
	case MSG_CMD_GET_WATCHDOG:
		return "WATCHDOG"
	case MSG_CMD_POE_RESET:
		return "POE_RESET"
	case MSG_CMD_SIM_TYPE:
		return "SIM_TYPE"
	case MSG_CMD_SIM_ICCID:
		return "SIM_ICCID"
	case MSG_CMD_SIM_CARRIER:
		return "SIM_CARRIER"
	case MSG_CMD_MODEM_SIGNAL:
		return "MODEM_SIGNAL"
	case MSG_CMD_MODEM_RESET:
		return "MODEM_RESET"
	case MSG_CMD_SERIAL_NUMBER:
		return "SERIAL_NUMBER"
	case MSG_CMD_BATCH_NUMBER:
		return "BATCH_NUMBER"
	case MSG_CMD_ANATEL_NUMBER:
		return "ANATEL_NUMBER"
	case MSG_CMD_OS_VERSION:
		return "OS_VERSION"
	case MSG_CMD_FIRMWARE_VERSION:
		return "FIRMWARE_VERSION"
	case MSG_CMD_BUZZER_ENABLE:
		return "BUZZER_ENABLE"
	case MSG_CMD_BUZZER_DISABLE:
		return "BUZZER_DISABLE"
	case MSG_CMD_IS_UPGRADING:
		return "IS_UPGRADING"
	case MSG_CMD_PCB_REV:
		return "PCB_REV"
	case MSG_CMD_TEST_START:
		return "TEST_START"
	case MSG_CMD_TEST_MODEM_RESULT:
		return "TEST_MODEM_RESULT"
	case MSG_CMD_TEST_SWITCH_RESULT:
		return "TEST_SWITCH_RESULT"
	case MSG_CMD_TEST_ETHERNET_RESULT:
		return "TEST_ETHERNET_RESULT"
	case MSG_CMD_HAS_BMS:
		return "HAS_BMS"
	case MSG_CMD_POWER_SOURCE:
		return "POWER_SOURCE"
	case MSG_CMD_SOCKETXP_STATUS:
		return "SOCKETXP_STATUS"
	case MSG_CMD_TAMPER_EVENT:
		return "TAMPER_EVENT"
	case MSG_CMD_BATTERY_LEVEL:
		return "BATTERY_LEVEL"
	case MSG_CMD_STM32_TEMPERATURE:
		return "STM32_TEMPERATURE"
	case MSG_CMD_TEST_EEPROM_RESULT:
		return "TEST_EEPROM_RESULT"
	case MSG_CMD_TEST_DISPLAY_RESULT:
		return "TEST_DISPLAY_RESULT"
	case MSG_CMD_EEPROM_DISABLE_WRITE_PROTECTION:
		return "MSG_CMD_EEPROM_DISABLE_WRITE_PROTECTION"
	case MSG_CMD_EEPROM_ENABLE_WRITE_PROTECTION:
		return "MSG_CMD_EEPROM_ENABLE_WRITE_PROTECTION"
	default:
		return "undefined"
	}
}
