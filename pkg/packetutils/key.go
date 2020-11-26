package packetutils

import (
	"encoding/hex"
	"strings"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	y3 "github.com/yomorun/yomo-codec-golang"
)

var (
	logger = utils.Logger.WithPrefix(utils.DefaultLogger, "yomoCodec")
)

func MatchingKey(key byte, node *y3.NodePacket) (flag bool, isNode bool, packet interface{}) {
	if len(node.PrimitivePackets) > 0 {
		for _, p := range node.PrimitivePackets {
			if key == p.SeqID() {
				return true, false, p
			}
		}
	}

	if len(node.NodePackets) > 0 {
		for _, n := range node.NodePackets {
			if key == n.SeqID() {
				return true, true, n
			}
			//return matchingKey(key, &n)
			flag, isNode, packet = MatchingKey(key, &n)
			if flag {
				return
			}
		}
	}

	return false, false, nil
}

func KeyOf(hexStr string) byte {
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
