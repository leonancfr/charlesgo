package main

import (
	"api"
	"charles_communicator"
	"device_info"
	"flag"
	"gablogger"
	"initializer"
	"monitor"
	mqtt "mqtt_connector"
	"peripherals"
	"socketxp"
	"sync"
	"updater"
	"utils"

	mqttPaho "github.com/eclipse/paho.mqtt.golang"
)

var Logger = gablogger.Logger()

func main() {
	//PARSE ARGUMENTS
	initFilePath := flag.String("config", "", "Specify the file path for initialization")
	flag.Parse()
	initializer.LoadConfig(*initFilePath)
	var mqtt_client_ptr *mqttPaho.Client = nil

	// MODULES INITIALIZATION
	deviceId, err := device_info.GetDeviceId()
	if err != nil {
		Logger.Fatal("Cannot retrieve MAC address")
	}
	gablogger.ConfigureDatadog(deviceId)

	// Initialize and get a pointer to a CharlesCommunicatorHandler instance
	charles_communicator.InitCC()
	CCHandler := charles_communicator.GetCharlesCommunicatorHandler()
	go CCHandler.Start()

	// Initialize API
	go api.Start()

	// Inject the CharlesCommunicatorHandler to other components
	peripherals.Setup(CCHandler)

	utils.Banner()
	mqtt.InitMQTTClient()
	updater.InitUpdater()
	socketxp.InitSocketXP()

	// Send OS version and Serial Number to STM32 display
	osVersion, _ := device_info.GetOSVersion()
	peripherals.SetOsVersion(osVersion)
	peripherals.SetSerialNumber(deviceId)

	// MQTT SETUP
	// register topics to subscribe
	mqtt_client_ptr = mqtt.GetMQTTClient()
	mqtt.RegisterSubscription(socketxp.Handler.CredentialTopic, socketxp.UpdateCredentialsCallback)
	mqtt.RegisterSubscription(updater.Handler.Hlk7628Topic, updater.UpdaterHlk7628Callback)
	mqtt.RegisterSubscription(updater.Handler.Stm32Topic, updater.UpdaterStm32Callback)
	// connect to broker after registering topics
	mqtt.Connect(mqtt_client_ptr)

	// OBSERVABILITY
	monitor.Monitor(mqtt_client_ptr)

	// INFINITE LOOP NOT TO EXIT MAIN
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
