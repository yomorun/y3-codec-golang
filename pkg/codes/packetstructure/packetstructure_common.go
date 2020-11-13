package packetstructure

import (
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

var (
	startingToken byte = 0x01
	logger             = utils.Logger.WithPrefix(utils.DefaultLogger, "yomoCodec::packetStructure")
)

type field struct {
	field reflect.StructField
	val   reflect.Value
}

func keyOf(hexStr string) byte {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = strings.TrimPrefix(hexStr, "0x")
	} else if strings.HasPrefix(hexStr, "0X") {
		hexStr = strings.TrimPrefix(hexStr, "0X")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		logger.Errorf("hex.DecodeString error: %v", err)
		return 0xff
	}

	if len(data) == 0 {
		logger.Errorf("hex.DecodeString data is []")
		return 0xff
	}

	return data[0]
}
