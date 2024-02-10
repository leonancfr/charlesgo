package utils

import (
	"encoding/hex"
	"testing"
)

func TestEncryptAES256ECB(t *testing.T) {
	key := []byte("12345612345612345612345612345612")
	data := []byte("11.22.33.44.55.66")

	expected_value := "e823d3ef05176c50cb852a6cac655590f85a69080ad64e19665d2208100dc47c"

	encrypted, err := EncryptAES256ECB(key, data)
	if err != nil {
		t.Fatalf("Encryption error: %v", err)
	}
	hexString := hex.EncodeToString(encrypted)

	if expected_value != hexString {
		t.Errorf("Expected decrypted data: %s, got: %s", string(expected_value), string(hexString))
	}
}
