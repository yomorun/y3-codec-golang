package y3

import (
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

type basePacket struct {
	length int64
	// Raw bytes
	raw []byte
}

func (bp *basePacket) Length() int64 {
	return bp.length
}

func (bp *basePacket) Buffer() []byte {
	return bp.raw
}

// isNodePacket determines if the packet is NodePacket or PrimitivePacket
func isNodePacket(flag byte) bool {
	return flag&utils.MSB == utils.MSB
}
