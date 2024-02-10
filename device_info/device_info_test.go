package device_info

import (
	"os"
	"testing"
)

func TestGetValueFromFile(t *testing.T) {
	tests := []struct {
		inputFileContent string
		key              string
		expectedValue    string
		expectError      bool
	}{
		{
			inputFileContent: `ID="ubuntu"
				ID_LIKE="debian"
				VERSION="20.04 LTS (Focal Fossa)"
				`,
			key:           "VERSION",
			expectedValue: "20.04 LTS (Focal Fossa)",
			expectError:   false,
		},
		{
			inputFileContent: `ID="centos"
				ID_LIKE="rhel fedora"
				VERSION="7 (Core)"
			`,
			key:           "ID",
			expectedValue: "centos",
			expectError:   false,
		},
		{
			inputFileContent: `ID="arch"
				ID_LIKE="unknown"
			`,
			key:           "VERSION",
			expectedValue: "",
			expectError:   true,
		},
		{
			inputFileContent: `invalid_line`,
			key:              "VERSION",
			expectedValue:    "",
			expectError:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			tmpFile, tmpFilePath := createTempFile(t, test.inputFileContent)
			defer tmpFile.Close()
			defer os.Remove(tmpFilePath)

			value, err := getValueFromFile(tmpFilePath, test.key)

			if test.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}

			if !test.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if value != test.expectedValue {
				t.Errorf("Expected value: %s, but got: %s", test.expectedValue, value)
			}
		})
	}
}

func createTempFile(t *testing.T, content string) (*os.File, string) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "os-release")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}

	// Write the simulated content to the temporary file
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	return tmpFile, tmpFile.Name()
}
