package y3

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

// IsNodePacket determines if the packet is NodePacket or PrimitivePacket
func IsNodePacket(flag byte) bool {
	return flag&MSB == MSB
}
