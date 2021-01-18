package utils

import (
	"encoding/hex"
	"strings"
)

// KeyOf: parse hex string to byte
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

func IsEmptyKey(observe byte) bool {
	return observe == byte(0)
}
