package spec

import (
	"errors"

	"github.com/yomorun/y3-codec-golang/pkg/encoding"
)

// FromBytes read Y3 buffer
func FromBytes(buf []byte) (p *Packet, err error) {
	if len(buf) < 2 {
		return nil, errors.New("malformed data")
	}

	p = &Packet{buffer: buf}
	pos := 0

	// read Tag Buffer, Tag support PVarUInt64 and raw bytes
	// if tagbuf is empty, then check idTag value
	cursor := readVariantLengthBuffer(buf, pos)
	p.tagbuf = buf[pos:cursor]
	pos = cursor

	// read Length
	var length uint64
	cursor, err = readPVarUInt64(buf, pos, &length)
	if err != nil {
		return nil, err
	}
	p.Length = length
	p.lenbuf = buf[pos:cursor]
	pos = cursor

	// read Value bytes
	cursor = pos + int(p.Length)
	if cursor > len(buf) {
		return nil, errors.New("malformed valbuf")
	}
	p.valbuf = buf[pos:cursor]

	return p, nil
}

func readVariantLengthBuffer(buffer []byte, position int) int {
	buf := buffer[position:]
	// PVarUInt64 type, MSB(0x80) is continuation bit
	cursor := 1
	for i, v := range buf {
		if v&0x80 != 0x80 {
			cursor += i
			break
		}
	}
	return cursor
}

func readPVarUInt64(buffer []byte, position int, val *uint64) (cursor int, err error) {
	buf := buffer[position:]
	// tag/length is PVarUInt64 type, MSB(0x80) is continuation bit
	cursor = 1
	for i, v := range buf {
		if v&0x80 != 0x80 {
			cursor += i
			break
		}
	}
	// generate tag buffer
	bytes := buf[:cursor]
	// read as PVarInt64
	codec := encoding.VarCodec{}
	err = codec.DecodePVarUInt64(bytes, val)
	return cursor + position, err
}

// GetValueAsUInt32 decode value as uint32
func (p *Packet) GetValueAsUInt32() (uint32, error) {
	var val uint32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt32(p.valbuf, &val)
	return val, err
}

// GetValueAsInt32 decode value as int32
func (p *Packet) GetValueAsInt32() (int32, error) {
	var val int32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt32(p.valbuf, &val)
	return val, err
}

// GetValueAsUInt64 decode value as uint64
func (p *Packet) GetValueAsUInt64() (uint64, error) {
	var val uint64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt64(p.valbuf, &val)
	return val, err
}

// GetValueAsInt64 decode value as int64
func (p *Packet) GetValueAsInt64() (int64, error) {
	var val int64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt64(p.valbuf, &val)
	return val, err
}

// GetValueAsFloat32 decode value as uint32
func (p *Packet) GetValueAsFloat32() (float32, error) {
	var val float32
	codec := encoding.VarCodec{Size: int(p.Length)}
	err := codec.DecodeVarFloat32(p.valbuf, &val)
	return val, err
}

// GetValueAsBool decode value as bool
func (p *Packet) GetValueAsBool() (bool, error) {
	res, err := p.GetValueAsUInt64()
	if res == 1 {
		return true, err
	}
	return false, err
}

// GetValueAsFloat64 decode value as float64
func (p *Packet) GetValueAsFloat64() (float64, error) {
	var val float64
	codec := encoding.VarCodec{Size: int(p.Length)}
	err := codec.DecodeVarFloat64(p.valbuf, &val)
	return val, err
}

// GetValueAsUTF8String decode value as float32
func (p *Packet) GetValueAsUTF8String() string {
	return string(p.valbuf)
}

// GetValueAsRawBytes decode value as float32
func (p *Packet) GetValueAsRawBytes() []byte {
	return p.valbuf
}
