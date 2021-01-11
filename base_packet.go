package y3

import (
	"github.com/yomorun/yomo-codec-golang/internal/mark"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

type basePacket struct {
	tag    *mark.Tag
	length uint32
	valbuf []byte
}

func (bp *basePacket) Length() uint32 {
	return bp.length
}

func (bp *basePacket) SeqID() byte {
	return bp.tag.SeqID()
}

// isNodePacket determines if the packet is NodePacket or PrimitivePacket
func isNodePacket(flag byte) bool {
	return flag&utils.MSB == utils.MSB
}
