package monitor

import (
	"device_info"
	"event_control"
	"gablogger"
	"network_info"
	"peripherals"
	"scheduler"
	"socketxp"
	"time"
	"updater"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Logger = gablogger.Logger()

var MQTTClient mqtt.Client

func Monitor(mqtt_client_ptr *mqtt.Client) {
	Logger.Infoln("Starting Monitor")
	MQTTClient = *mqtt_client_ptr

	event_control.RegisterToReceiveEvent(peripherals.GetTamperEventId(), sendTamperEvent, nil)
	event_control.RegisterToReceiveEvent(peripherals.GetPowerSourceEventId(), sendPowerSourceEvent, nil)
	event_control.RegisterToReceiveEvent(peripherals.GetBuzzerEventId(), sendBuzzerEvent, nil)

	event_control.RegisterToReceiveEvent(updater.GetSTM32pdateEventId(), sendUpdateSTM32Event, nil)
	event_control.RegisterToReceiveEvent(updater.GetHLK7628UpdateEventId(), sendUpdateHLK7628Event, nil)

	publishMetricFromFunction(topicOsVersion, device_info.GetOSVersion)
	publishMetricFromFunction(topicStm32FirmwareVersion, peripherals.GetFirmwareVersion)
	publishMetricFromFunction(topicCharlesGoVersion, device_info.GetCharlesGoVersion)

	if scheduler.InitScheduler() != nil {
		Logger.Error("Cannot initialize scheduler")
		return
	}

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicSocketxpStatus, socketxp.IsConnected)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicMacAddress, device_info.GetDeviceId)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicOsVersion, device_info.GetOSVersion)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicSTM32Temperature, peripherals.GetSTM32Temperature)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicStm32FirmwareVersion, peripherals.GetFirmwareVersion)

	// Delay for 3 seconds using time.Sleep to prevent sending multiple messages simultaneously to the STM32, causing message timeout
	// This ensures a pause between messages and avoids potential issues.
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicPcbBatchNumber, peripherals.GetPCBBatch)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicPcbReview, peripherals.GetPCBReview)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicSimCardCarrier, peripherals.GetSIMCardCarrier)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicHasBMS, peripherals.GetHasBMS)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicBatteryLevel, peripherals.GetBatteryLevel)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicIccid, peripherals.GetSIMCardICCID)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicModemSignal, peripherals.GetModemSignalStrength)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicSimCardType, peripherals.GetSIMCardType)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicModemConnType, peripherals.GetModemConnectionType)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicModemConnBand, peripherals.GetModemConnectionBand)
	time.Sleep(time.Second * 3)

	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicConnectionInUse, network_info.GetPriorityRoute)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicModemConnectionStatus, network_info.GetModemInterfaceStatus)
	scheduler.RegisterFunctionToSchedule(time.Minute*5, publishMetricFromFunction, topicWiredonnectionStatus, network_info.GetWiredInterfaceStatus)

}
