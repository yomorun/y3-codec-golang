package encoding

import (
	"errors"
)

func SizeOfNVarInt32(value int32) int {
	return sizeOfNVarInt(int64(value), 32)
}

func (codec *VarIntCodec) EncodeNVarInt32(buffer []byte, value int32) error {
	return codec.encodeNVarInt(buffer, int64(value))
}

func (codec *VarIntCodec) DecodeNVarInt32(buffer []byte, value *int32) error {
	var val = int64(*value)
	var err = codec.decodeNVarInt(buffer, &val)
	*value = int32(val)
	return err
}

func SizeOfNVarUInt32(value uint32) int {
	return sizeOfNVarInt(int64(int32(value)), 32)
}

func (codec *VarIntCodec) EncodeNVarUInt32(buffer []byte, value uint32) error {
	return codec.encodeNVarInt(buffer, int64(int32(value)))
}

func (codec *VarIntCodec) DecodeNVarUInt32(buffer []byte, value *uint32) error {
	var val = int64(int32(*value))
	var err = codec.decodeNVarInt(buffer, &val)
	*value = uint32(val)
	return err
}

func SizeOfNVarInt64(value int64) int {
	return sizeOfNVarInt(value, 64)
}

func (codec *VarIntCodec) EncodeNVarInt64(buffer []byte, value int64) error {
	return codec.encodeNVarInt(buffer, value)
}

func (codec *VarIntCodec) DecodeNVarInt64(buffer []byte, value *int64) error {
	return codec.decodeNVarInt(buffer, value)
}

func SizeOfNVarUInt64(value uint64) int {
	return sizeOfNVarInt(int64(value), 64)
}

func (codec *VarIntCodec) EncodeNVarUInt64(buffer []byte, value uint64) error {
	return codec.encodeNVarInt(buffer, int64(value))
}

func (codec *VarIntCodec) DecodeNVarUInt64(buffer []byte, value *uint64) error {
	var val = int64(*value)
	var err = codec.decodeNVarInt(buffer, &val)
	*value = uint64(val)
	return err
}

func sizeOfNVarInt(value int64, width int) int {
	const unit = 8 // 编码组位宽
	var lead = value >> (width - 1)

	for size := width / unit; size > 0; size-- {
		var lookAhead = value >> (size*unit - 1)
		if lookAhead != lead {
			return size + 1
		}
	}
	return 1
}

func (codec *VarIntCodec) encodeNVarInt(buffer []byte, value int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}

	const unit = 8
	for codec.Size > 0 {
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}

		codec.Size--
		buffer[codec.Ptr] = byte(value >> (codec.Size * unit))
		codec.Ptr++
	}
	return nil
}

func (codec *VarIntCodec) decodeNVarInt(buffer []byte, value *int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to decode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 8
	if codec.Size > 0 { // 初始化符号
		*value = int64(int8(buffer[codec.Ptr]) >> 7)
		codec.Size = -codec.Size
	}
	for codec.Size < 0 {
		codec.Size++
		*value = (*value << unit) | int64(buffer[codec.Ptr])
		codec.Ptr++
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
	return nil
}
