package initializer

type config struct {
	deviceConfig deviceConfig
	supervisor   supervisorConfig
	mqtt         mqttConfig
	updater      updaterConfig
}

type supervisorConfig struct {
	IsEnabled bool
}

type deviceConfig struct {
	Label string
}

type mqttConfig struct {
	IsEnabled bool
}

type updaterConfig struct {
	IsEnabledStm32   bool
	IsEnabledHlk7628 bool
}
