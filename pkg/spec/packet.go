package spec

// Packet is a TLV group
type Packet struct {
	Length uint64
	tagbuf []byte
	lenbuf []byte
	valbuf []byte
	buffer []byte
	idTag  uint64
}

// GetTag return Tag as uint64 value
func (p *Packet) GetTag() uint64 {
	var tag uint64
	readPVarUInt64(p.tagbuf, 0, &tag)
	return tag
}

// GetTagBuffer return Tag as raw bytes
func (p *Packet) GetTagBuffer() []byte {
	return p.tagbuf
}

// GetLengthBuffer return Tag as raw bytes
func (p *Packet) GetLengthBuffer() []byte {
	return p.lenbuf
}

// GetValueBuffer return Tag as raw bytes
func (p *Packet) GetValueBuffer() []byte {
	return p.valbuf
}
