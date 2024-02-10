package common

import (
	"os"
	"strconv"
)

var (
	VERSION         string
	ENVIRONMENT     string
	MAGIC_KEY       string
	MQTT_BROKER     string
	MQTT_PORT       string
	SFTP_SERVER     string
	SFTP_PORT       string
	SFTP_USER       string
	SFTP_PASS       string
	DATADOG_HOST    string
	DATADOG_API_KEY string
	API_PORT        string
)

func init() {
	initializeVar(&VERSION, "VERSION")
	initializeVar(&ENVIRONMENT, "ENVIRONMENT")
	initializeVar(&MAGIC_KEY, "MAGIC_KEY")
	initializeVar(&MQTT_BROKER, "MQTT_BROKER")
	initializeVar(&MQTT_PORT, "MQTT_PORT")
	initializeVar(&SFTP_SERVER, "SFTP_SERVER")
	initializeVar(&SFTP_PORT, "SFTP_PORT")
	initializeVar(&SFTP_USER, "SFTP_USER")
	initializeVar(&SFTP_PASS, "SFTP_PASS")
	initializeVar(&DATADOG_HOST, "DATADOG_HOST")
	initializeVar(&DATADOG_API_KEY, "DATADOG_API_KEY")
	initializeVar(&API_PORT, "API_PORT")
}

func initializeVar(constVar *string, envName string) {
	if *constVar == "" {
		*constVar = os.Getenv(envName)
	}
}

func initializeVarInt(constVar *int, envName string) {
	if *constVar == 0 {
		valStr := os.Getenv(envName)
		if val, err := strconv.Atoi(valStr); err == nil {
			*constVar = val
		}
	}
}
