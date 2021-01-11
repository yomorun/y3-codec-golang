package codes

import (
	"errors"
	"fmt"

	"github.com/yomorun/y3-codec-golang/pkg/packetutils"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/y3-codec-golang/pkg/codes/packetstructure"
)

// StructDecoder: for UnmarshalStruct
type StructDecoder struct {
	Observe byte
}

func newStructDecoder(observe byte) *StructDecoder {
	return &StructDecoder{Observe: observe}
}

func (d StructDecoder) Unmarshal(data []byte, mold interface{}) error {
	nodePacket, _, err := y3.DecodeNodePacket(data)
	if err != nil {
		return err
	}
	return d.UnmarshalByNodePacket(nodePacket, mold)
}

func (d StructDecoder) UnmarshalNative(data []byte, mold interface{}) error {
	nodePacket, _, err := y3.DecodeNodePacket(data)
	if err != nil {
		return err
	}

	err = packetstructure.Decode(nodePacket, mold)
	if err != nil {
		return err
	}
	return nil
}

func (d StructDecoder) UnmarshalByNodePacket(node *y3.NodePacket, mold interface{}) error {
	key := d.Observe
	ok, isNode, packet := packetutils.MatchingKey(key, node)
	if !ok {
		return errors.New(fmt.Sprintf("not found mold in result. key:%#x", key))
	}
	if !isNode {
		return errors.New(fmt.Sprintf("packet must be NodePacket. key:%#x", key))
	}

	nodePacket := packet.(y3.NodePacket)
	return packetstructure.Decode(&nodePacket, mold)
}
