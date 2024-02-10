package monitor

import (
	"encoding/json"
	"fmt"
	"time"
	"utils"
)

type mqttPublishMessage struct {
	Data string    `json:"data"`
	Date time.Time `json:"date"`
}

func publishMetricFromFunction(topic string, getMetricFunction func() (string, error)) {
	value, err := getMetricFunction()
	if err != nil {
		Logger.Error("Cannot get data to publish in topic ", topic, ". Reason: ", err)
		return
	}
	publishMetric(topic, value)
}

func publishMetric(topic, value string) {
	path := fmt.Sprintf("devices/%s/monitoring/%s", utils.GetUserName(), topic)
	payload := mqttPublishMessage{
		Data: value,
		Date: time.Now(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		Logger.Error("Error marshalling the MQTT payload: ", err)
		return
	}

	MQTTClient.Publish(path, 1, false, jsonPayload)
	Logger.Debug("Publish \"", value, "\" -> \"", path, "\" topic.")
}
