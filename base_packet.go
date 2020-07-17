package y3

import (
	codec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

type basePacket struct {
	tag    *codec.Tag
	length uint64
	valbuf []byte
}

func (bp *basePacket) Length() uint64 {
	return bp.length
}

// func (bp *basePacket) Buffer() []byte {
// 	return bp.valbuf
// }

func (bp *basePacket) SeqID() byte {
	return bp.tag.SeqID()
}

// isNodePacket determines if the packet is NodePacket or PrimitivePacket
func isNodePacket(flag byte) bool {
	return flag&utils.MSB == utils.MSB
}
