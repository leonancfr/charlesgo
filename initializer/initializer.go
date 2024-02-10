package initializer

import (
	"fmt"
	"gablogger"

	goIni "gopkg.in/ini.v1"
)

var ini config
var Logger = gablogger.Logger()

func LoadConfig(filePath string) {
	if filePath != "" {
		Logger.Debugln("Using", filePath, "as initialize file")
		cfg, err := goIni.Load(filePath)
		if err != nil {
			Logger.Fatalln("Cannot load", filePath, "file. Reason:", err)
		}
		loadSupervisorConfig(cfg)
		loadDeviceConfig(cfg)
		loadMqttConfig(cfg)
		loadUpdaterConfig(cfg)
	} else {
		initializeDefaultConfig()
	}
}

func initializeDefaultConfig() {
	ini.supervisor.IsEnabled = true
	ini.mqtt.IsEnabled = true
	ini.deviceConfig.Label = ""
	ini.updater.IsEnabledHlk7628 = true
	ini.updater.IsEnabledStm32 = true
}

func loadDeviceConfig(cfg *goIni.File) {
	ini.deviceConfig.Label = cfg.Section("DEVICE_INFO").Key("LABEL").String()
}

func loadUpdaterConfig(cfg *goIni.File) {
	var err error
	ini.updater.IsEnabledStm32, err = getBoolValue(cfg, "UPDATE", "ENABLE_STM32", true)
	if err != nil {
		Logger.WithField("invalid-value", "config-file").Errorln(err, "Using default value.")
	}

	ini.updater.IsEnabledHlk7628, err = getBoolValue(cfg, "UPDATE", "ENABLE_HLK7628", true)
	if err != nil {
		Logger.WithField("invalid-value", "config-file").Errorln(err, "Using default value.")
	}
}

func loadMqttConfig(cfg *goIni.File) {
	var err error
	ini.mqtt.IsEnabled, err = getBoolValue(cfg, "MQTT", "ENABLE", true)
	if err != nil {
		Logger.WithField("invalid-value", "config-file").Errorln(err, "Using default value.")
	}
}

func loadSupervisorConfig(cfg *goIni.File) {
	var err error
	ini.supervisor.IsEnabled, err = getBoolValue(cfg, "SUPERVISOR", "ENABLE", true)
	if err != nil {
		Logger.WithField("invalid-value", "config-file").Errorln(err, "Using default value.")
	}
}

func getBoolValue(cfg *goIni.File, section, key string, defaultValue bool) (bool, error) {
	rawEnableField := cfg.Section(section).Key(key)
	if rawEnableField != nil {
		isEnableField, err := rawEnableField.Bool()
		if err != nil {
			return defaultValue, fmt.Errorf("Cannot decode '%s %s'. %v", section, key, err)
		} else {
			return isEnableField, nil
		}
	}
	return defaultValue, fmt.Errorf("Field '%s %s' not found", section, key)
}

func GetLabel() string {
	return ini.deviceConfig.Label
}

func IsSupervisorEnable() bool {
	return ini.supervisor.IsEnabled
}

func IsMqttEnabled() bool {
	return ini.mqtt.IsEnabled
}

func IsHlk7628UpdateEnabled() bool {
	return ini.updater.IsEnabledHlk7628
}

func IsStm32UpdateEnabled() bool {
	return ini.updater.IsEnabledStm32
}
