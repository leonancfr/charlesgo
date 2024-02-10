package mqtt_connector

import (
	"common"
	"crypto/tls"
	"device_info"
	"encoding/hex"
	"fmt"
	"gablogger"
	"time"
	"utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var Logger = gablogger.Logger()
var client *mqtt.Client
var subscriptions = make(map[string]mqtt.MessageHandler)

func GetMQTTClient() *mqtt.Client {
	return client
}

// Function to set up and return a pointer to MQTT_Connector
func InitMQTTClient() {
	Logger.Debugln("mqtt Setup")

	var LabelMAC string
	var err error
	for retry := 0; retry < 3; retry++ {
		LabelMAC, err = device_info.GetDeviceId()
		if err == nil {
			break
		}
		if retry == 2 {
			Logger.Error("Cannot get MAC from device")
			panic(err)
		}
		time.Sleep(5)
	}
	bin_password, _ := utils.EncryptAES256ECB([]byte(common.MAGIC_KEY), []byte(LabelMAC))

	password := hex.EncodeToString(bin_password)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	username := utils.GetUserName()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", common.MQTT_BROKER, common.MQTT_PORT))
	opts.SetClientID(username)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetTLSConfig(tlsConfig)
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	Logger.Debugln("SetConnectionLostHandler")
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		Logger.Errorln("mqtt connection lost error: " + err.Error())
	}
	Logger.Debugln("SetReconnectingHandler")
	opts.OnReconnecting = func(c mqtt.Client, options *mqtt.ClientOptions) {
		Logger.Debugln("mqtt reconnecting")
	}
	Logger.Debugln("SetOnConnectHandler")
	opts.OnConnect = func(c mqtt.Client) {
		Logger.Infoln("mqtt connected to broker", common.MQTT_BROKER, common.MQTT_PORT)
		SubscribeAll()
	}

	c := mqtt.NewClient(opts)
	client = &c
}

func Connect(client *mqtt.Client) {
	for {
		if token := (*client).Connect(); token.Wait() && token.Error() != nil {
			Logger.Errorln("Could not connect", token.Error())
			time.Sleep(5 * time.Second)
		} else {
			Logger.Debugln("Connected")
			break
		}
	}
}

func RegisterSubscription(topic string, callback mqtt.MessageHandler) {
	subscriptions[topic] = callback
}

func SubscribeAll() {
	Logger.Infoln("Subscribing to registered topics")
	for topic, callback := range subscriptions {
		token := (*client).Subscribe(topic, 1, callback)
		token.Wait()
	}
}
