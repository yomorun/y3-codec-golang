package codes

import (
	"errors"
	"fmt"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	y3 "github.com/yomorun/yomo-codec-golang"
	"github.com/yomorun/yomo-codec-golang/pkg/codes/packetstructure"
)

// StructDecoder: for UnmarshalStruct
type StructDecoder struct {
	Observe string
}

func newStructDecoder(observe string) *StructDecoder {
	return &StructDecoder{Observe: observe}
}

// Unmarshal: Unmarshal []byte to interface
func (d StructDecoder) Unmarshal(data []byte, mold interface{}) error {
	key := packetutils.KeyOf(d.Observe)
	pct, _, err := y3.DecodeNodePacket(data)
	if err != nil {
		return err
	}

	ok, isNode, packet := packetutils.MatchingKey(key, pct)
	if !ok {
		return errors.New(fmt.Sprintf("not found mold in result. key:%#x", key))
	}
	if !isNode {
		return errors.New(fmt.Sprintf("packet must be NodePacket. key:%#x", key))
	}

	nodePacket := packet.(y3.NodePacket)
	return packetstructure.Decode(&nodePacket, mold)
}
