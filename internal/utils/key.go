package utils

import (
	"encoding/hex"
	"strings"
)

const (
	// EmptyKey mark an empty key
	EmptyKey byte = 0
)

// IsEmptyKey determine if observe is empty
func IsEmptyKey(observe byte) bool {
	return observe == byte(EmptyKey)
}

// KeyOf parse hex string to byte
func KeyOf(hexStr string) byte {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = strings.TrimPrefix(hexStr, "0x")
	} else if strings.HasPrefix(hexStr, "0X") {
		hexStr = strings.TrimPrefix(hexStr, "0X")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		DefaultLogger.Errorf("hex.DecodeString error: %v", err)
		return 0x00
	}

	if len(data) == 0 {
		DefaultLogger.Errorf("hex.DecodeString data is []")
		return 0x00
	}

	return data[0]
}

// ForbidUserKey forbid user set that key
func ForbidUserKey(key byte) bool {
	switch key {
	case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return true
	}
	return false
}

// AllowSignalKey allow set that signal key
func AllowSignalKey(key byte) bool {
	switch key {
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return true
	}
	return false
}
