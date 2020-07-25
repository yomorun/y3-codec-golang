package encoding

import (
	"errors"
)

func SizeOfPVarInt32(value int32) int {
	return sizeOfPVarInt(int64(value), 32)
}

func (codec *VarIntCodec) EncodePVarInt32(buffer []byte, value int32) error {
	return codec.encodePVarInt(buffer, int64(value))
}

func (codec *VarIntCodec) DecodePVarInt32(buffer []byte, value *int32) error {
	var val = int64(*value)
	var err = codec.decodePVarInt(buffer, &val)
	*value = int32(val)
	return err
}

func SizeOfPVarUInt32(value uint32) int {
	return sizeOfPVarInt(int64(int32(value)), 32)
}

func (codec *VarIntCodec) EncodePVarUInt32(buffer []byte, value uint32) error {
	return codec.encodePVarInt(buffer, int64(int32(value)))
}

func (codec *VarIntCodec) DecodePVarUInt32(buffer []byte, value *uint32) error {
	var val = int64(int32(*value))
	var err = codec.decodePVarInt(buffer, &val)
	*value = uint32(val)
	return err
}

func SizeOfPVarInt64(value int64) int {
	return sizeOfPVarInt(value, 64)
}

func (codec *VarIntCodec) EncodePVarInt64(buffer []byte, value int64) error {
	return codec.encodePVarInt(buffer, value)
}

func (codec *VarIntCodec) DecodePVarInt64(buffer []byte, value *int64) error {
	return codec.decodePVarInt(buffer, value)
}

func SizeOfPVarUInt64(value uint64) int {
	return sizeOfPVarInt(int64(value), 64)
}

func (codec *VarIntCodec) EncodePVarUInt64(buffer []byte, value uint64) error {
	return codec.encodePVarInt(buffer, int64(value))
}

func (codec *VarIntCodec) DecodePVarUInt64(buffer []byte, value *uint64) error {
	var val = int64(*value)
	var err = codec.decodePVarInt(buffer, &val)
	*value = uint64(val)
	return err
}

func sizeOfPVarInt(value int64, width int) int {
	const unit = 7 // 编码组位宽
	var lead = value >> (width - 1)

	for size := width / unit; size > 0; size-- {
		var lookAhead = value >> (size*unit - 1)
		if lookAhead != lead {
			return size + 1
		}
	}
	return 1
}

func (codec *VarIntCodec) encodePVarInt(buffer []byte, value int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 7
	const more = -1 << unit

	for codec.Size > 1 {
		codec.Size--
		var part = value >> (codec.Size * unit)

		buffer[codec.Ptr] = byte(part | more)
		codec.Ptr++

		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}

	const mask = -1 ^ more
	codec.Size = 0
	buffer[codec.Ptr] = byte(value & mask)
	codec.Ptr++
	return nil
}

func (codec *VarIntCodec) decodePVarInt(buffer []byte, value *int64) error {
	if codec == nil {
		return errors.New("nothing to decode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 7
	const mask = -1 ^ (-1 << unit)

	if codec.Size == 0 { // 初始化符号
		*value = int64(int8(buffer[codec.Ptr]) << (8 - unit) >> unit)
	}

	for {
		var part = int8(buffer[codec.Ptr])
		codec.Ptr++

		codec.Size++
		*value = (*value << unit) | int64(mask&part)

		if part >= 0 { // 最后一个字节
			return nil
		}
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
}
