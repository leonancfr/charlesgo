package socketxp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeMessageSuccess(t *testing.T) {
	expected_value_DeviceId := "device_id"
	expected_value_DeviceKey := "device_key"
	expected_value_DeviceName := "device_name"
	expected_value_DeviceGroup := ""

	message := `{"DeviceId":"device_id","DeviceKey":"device_key","DeviceName":"device_name"}`

	key, err := decodeMessage(message)

	if err != nil {
		t.Errorf("Expected no erros, got: %v", err)
	}

	if key.DeviceName != expected_value_DeviceName {
		t.Errorf("Expected data in DeviceName: %s, got: %s", string(expected_value_DeviceName), string(key.DeviceName))
	}

	if key.DeviceId != expected_value_DeviceId {
		t.Errorf("Expected data in DeviceId: %s, got: %s", string(expected_value_DeviceId), string(key.DeviceId))
	}

	if key.DeviceKey != expected_value_DeviceKey {
		t.Errorf("Expected data in DeviceKey: %s, got: %s", string(expected_value_DeviceKey), string(key.DeviceKey))
	}

	if key.DeviceGroup != expected_value_DeviceGroup {
		t.Errorf("Expected data in DeviceGroup: %s, got: %s", string(expected_value_DeviceGroup), string(key.DeviceGroup))
	}
}

func TestDecodeMessageErrors(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		expected bool // true se um erro for esperado, false caso contrário
	}{
		{
			name:     "MissingDeviceId",
			message:  `{"DeviceKey":"device_key","DeviceName":"device_name"}`,
			expected: true,
		},
		{
			name:     "MissingDeviceKey",
			message:  `{"DeviceId":"device_id","DeviceName":"device_name"}`,
			expected: true,
		},
		{
			name:     "MissingDeviceName",
			message:  `{"DeviceId":"device_id","DeviceKey":"device_key"}`,
			expected: true,
		},
		{
			name:     "MissingAllFields",
			message:  `{}`,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decodeMessage(tc.message)

			if tc.expected && err == nil {
				t.Error("Expected an error, but got nil")
			} else if !tc.expected && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

func TestDeviceKeyIsEqual(t *testing.T) {
	device1 := deviceKey{
		DeviceId:    "device_id_1",
		DeviceKey:   "device_key_1",
		DeviceName:  "device_name_1",
		DeviceGroup: "group_A",
	}

	device2 := deviceKey{
		DeviceId:    "device_id_1",
		DeviceKey:   "device_key_1",
		DeviceName:  "device_name_1",
		DeviceGroup: "group_A",
	}

	device3 := deviceKey{
		DeviceId:    "device_id_2",
		DeviceKey:   "device_key_2",
		DeviceName:  "device_name_2",
		DeviceGroup: "group_B",
	}

	t.Run("EqualDevices", func(t *testing.T) {
		if !device1.isEqual(device2) {
			t.Error("Expected two equal devices to be considered equal.")
		}
	})

	t.Run("DifferentDevices", func(t *testing.T) {
		if device1.isEqual(device3) {
			t.Error("Expected two different devices to be considered not equal.")
		}
	})
}

func TestConfigJsonIsEqual(t *testing.T) {
	config1 := configJson{
		WorkDir: "work_dir_1",
		Tunnels: []tunnel{
			{Destination: "dest_1", CustomDomain: "custom_1", Subdomain: "sub_1"},
			{Destination: "dest_2", CustomDomain: "custom_2", Subdomain: "sub_2"},
		},
	}

	config2 := configJson{
		WorkDir: "work_dir_1",
		Tunnels: []tunnel{
			{Destination: "dest_1", CustomDomain: "custom_1", Subdomain: "sub_1"},
			{Destination: "dest_2", CustomDomain: "custom_2", Subdomain: "sub_2"},
		},
	}

	config3 := configJson{
		WorkDir: "work_dir_2",
		Tunnels: []tunnel{
			{Destination: "dest_1", CustomDomain: "custom_1", Subdomain: "sub_1"},
			{Destination: "dest_2", CustomDomain: "custom_2", Subdomain: "sub_2"},
		},
	}

	t.Run("EqualConfigs", func(t *testing.T) {
		if !config1.isEqual(config2) {
			t.Error("Expected two equal configs to be considered equal.")
		}
	})

	t.Run("DifferentConfigs", func(t *testing.T) {
		if config1.isEqual(config3) {
			t.Error("Expected two different configs to be considered not equal.")
		}
	})
}

func TestLoadDeviceKeyFile(t *testing.T) {
	// Criar um diretório temporário para testar
	tmpDir := t.TempDir()

	// Criar o arquivo deviceKey no diretório temporário
	deviceKeyData := `{"DeviceId":"TEST_DEVICEID","DeviceKey":"TEST_DEVICE_KEY","DeviceName":"TEST_DEVICE_NAME","DeviceGroup":"TEST_DEVICE_GROUP"}`
	deviceKeyPath := filepath.Join(tmpDir, "device.key")
	err := os.WriteFile(deviceKeyPath, []byte(deviceKeyData), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar o arquivo de chave do dispositivo: %v", err)
	}

	socketXp := SocketXp{
		deviceKeyPath: deviceKeyPath,
	}

	loadedDeviceKey, err := socketXp.loadDeviceKeyFile()
	if err != nil {
		t.Fatalf("Erro ao carregar a chave do dispositivo: %v", err)
	}
	fmt.Println(loadedDeviceKey)
	expectedDeviceKey := deviceKey{
		DeviceId:    "TEST_DEVICEID",
		DeviceKey:   "TEST_DEVICE_KEY",
		DeviceName:  "TEST_DEVICE_NAME",
		DeviceGroup: "TEST_DEVICE_GROUP",
	}

	if loadedDeviceKey != expectedDeviceKey {
		t.Errorf("Chave do dispositivo carregada incorreta. Esperada: %v, Recebida: %v", expectedDeviceKey, loadedDeviceKey)
	}
}

func TestLoadDeviceKeyFileNonExistent(t *testing.T) {
	// Criar um diretório temporário para testar
	tmpDir := t.TempDir()

	socketXp := SocketXp{
		deviceKeyPath: filepath.Join(tmpDir, "non_existent.json"),
	}

	// Testar a função loadDeviceKeyFile com um arquivo inexistente
	_, err := socketXp.loadDeviceKeyFile()
	if err == nil {
		t.Error("A função não retornou erro ao carregar um arquivo inexistente, mas deveria.")
	}
}

func TestLoadConfigJsonFile(t *testing.T) {
	// Criar um diretório temporário para testar
	tmpDir := t.TempDir()

	// Criar o arquivo configJson no diretório temporário
	configData := `{
		"work_dir": "TEST_WORKDIR",
		"tunnels": [
			{
				"destination": "TEST_TUNNEL_DESTINATION"
			},
			{
				"destination": "TEST_DESTINATION",
				"custom_domain": "TEST_CUSTOM_DOMAIN",
				"subdomain": "TEST_SUBDOMAIN"
			}
		]
	}`

	configJSONPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configJSONPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar o arquivo de configuração JSON: %v", err)
	}

	socketXp := SocketXp{
		configJSONPath: configJSONPath,
	}

	// Testar a função loadConfigJsonFile
	config, err := socketXp.loadConfigJsonFile()
	if err != nil {
		t.Fatalf("Erro ao carregar a configuração JSON: %v", err)
	}

	expectedConfig := configJson{
		WorkDir: "TEST_WORKDIR",
		Tunnels: []tunnel{
			{
				Destination: "TEST_TUNNEL_DESTINATION",
			},
			{
				Destination:  "TEST_DESTINATION",
				CustomDomain: "TEST_CUSTOM_DOMAIN",
				Subdomain:    "TEST_SUBDOMAIN",
			},
		},
	}

	if !config.isEqual(expectedConfig) {
		t.Errorf("Configuração carregada incorreta. Esperada: %v, Recebida: %v", expectedConfig, config)
	}
}

func TestLoadConfigJsonFileNonExistent(t *testing.T) {
	// Criar um diretório temporário para testar
	tmpDir := t.TempDir()

	socketXp := SocketXp{
		configJSONPath: filepath.Join(tmpDir, "non_existent.json"),
	}

	// Testar a função loadConfigJsonFile com um arquivo inexistente
	_, err := socketXp.loadConfigJsonFile()
	if err == nil {
		t.Error("A função não retornou erro ao carregar um arquivo inexistente, mas deveria.")
	}
}
