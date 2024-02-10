package charles_communicator

const (
	BUFFER_SIZE         = 100
	RECEIVING_DATA      = 0
	WAITING_FOR_MESSAGE = 1
	PROTOCOL_VERSION    = 0
)

const (
	MSG_TYPE_GET     = 0
	MSG_TYPE_SET     = 1
	MSG_TYPE_RESP    = 2
	MSG_TYPE_ERROR   = 3
	MSG_TYPE_TIMEOUT = 4
)

const (
	MSG_CMD_GET_WATCHDOG                    = 0
	MSG_CMD_POE_RESET                       = 1
	MSG_CMD_SIM_TYPE                        = 2
	MSG_CMD_SIM_ICCID                       = 3
	MSG_CMD_SIM_CARRIER                     = 4
	MSG_CMD_MODEM_SIGNAL                    = 5
	MSG_CMD_MODEM_RESET                     = 6
	MSG_CMD_SERIAL_NUMBER                   = 7
	MSG_CMD_BATCH_NUMBER                    = 8
	MSG_CMD_ANATEL_NUMBER                   = 9
	MSG_CMD_OS_VERSION                      = 10
	MSG_CMD_FIRMWARE_VERSION                = 11
	MSG_CMD_BUZZER_ENABLE                   = 12
	MSG_CMD_BUZZER_DISABLE                  = 13
	MSG_CMD_IS_UPGRADING                    = 14
	MSG_CMD_PCB_REV                         = 15
	MSG_CMD_TEST_START                      = 16
	MSG_CMD_TEST_MODEM_RESULT               = 17
	MSG_CMD_TEST_SWITCH_RESULT              = 18
	MSG_CMD_TEST_ETHERNET_RESULT            = 19
	MSG_CMD_HAS_BMS                         = 20
	MSG_CMD_POWER_SOURCE                    = 21
	MSG_CMD_SOCKETXP_STATUS                 = 22
	MSG_CMD_TAMPER_EVENT                    = 23
	MSG_CMD_TEST_EEPROM_RESULT              = 24
	MSG_CMD_TEST_DISPLAY_RESULT             = 25
	MSG_CMD_EEPROM_DISABLE_WRITE_PROTECTION = 26
	MSG_CMD_EEPROM_ENABLE_WRITE_PROTECTION  = 27
	MSG_CMD_STM32_TEMPERATURE               = 28
	MSG_CMD_BATTERY_LEVEL                   = 29
	MSG_CMD_MODEM_CONN_TYPE                 = 30
	MSG_CMD_MODEM_CONN_BAND                 = 31
)

const WAIT_MESSAGE_RESPONSE_TIMEOUT = 10000
