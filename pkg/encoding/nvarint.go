package encoding

import (
	"errors"
)

// SizeOfNVarInt32 return the buffer size after encoding value as NVarInt32
func SizeOfNVarInt32(value int32) int {
	return sizeOfNVarInt(int64(value), 32)
}

// EncodeNVarInt32 encode value as NVarInt32 to buffer
func (codec *VarCodec) EncodeNVarInt32(buffer []byte, value int32) error {
	return codec.encodeNVarInt(buffer, int64(value))
}

// DecodeNVarInt32 decode to value as NVarInt32 from buffer
func (codec *VarCodec) DecodeNVarInt32(buffer []byte, value *int32) error {
	var val = int64(*value)
	var err = codec.decodeNVarInt(buffer, &val)
	*value = int32(val)
	return err
}

// SizeOfNVarUInt32 return the buffer size after encoding value as NVarUInt32
func SizeOfNVarUInt32(value uint32) int {
	return sizeOfNVarInt(int64(int32(value)), 32)
}

// EncodeNVarUInt32 encode value as NVarUInt32 to buffer
func (codec *VarCodec) EncodeNVarUInt32(buffer []byte, value uint32) error {
	return codec.encodeNVarInt(buffer, int64(int32(value)))
}

// DecodeNVarUInt32 decode to value as NVarUInt32 from buffer
func (codec *VarCodec) DecodeNVarUInt32(buffer []byte, value *uint32) error {
	var val = int64(int32(*value))
	var err = codec.decodeNVarInt(buffer, &val)
	*value = uint32(val)
	return err
}

// SizeOfNVarInt64 return the buffer size after encoding value as NVarInt64
func SizeOfNVarInt64(value int64) int {
	return sizeOfNVarInt(value, 64)
}

// EncodeNVarInt64 encode value as NVarInt64 to buffer
func (codec *VarCodec) EncodeNVarInt64(buffer []byte, value int64) error {
	return codec.encodeNVarInt(buffer, value)
}

// DecodeNVarInt64 decode to value as NVarInt64 from buffer
func (codec *VarCodec) DecodeNVarInt64(buffer []byte, value *int64) error {
	return codec.decodeNVarInt(buffer, value)
}

// SizeOfNVarUInt64 return the buffer size after encoding value as NVarUInt64
func SizeOfNVarUInt64(value uint64) int {
	return sizeOfNVarInt(int64(value), 64)
}

// EncodeNVarUInt64 encode value as NVarUInt64 to buffer
func (codec *VarCodec) EncodeNVarUInt64(buffer []byte, value uint64) error {
	return codec.encodeNVarInt(buffer, int64(value))
}

// DecodeNVarUInt64 decode to value as NVarUInt64 from buffer
func (codec *VarCodec) DecodeNVarUInt64(buffer []byte, value *uint64) error {
	var val = int64(*value)
	var err = codec.decodeNVarInt(buffer, &val)
	*value = uint64(val)
	return err
}

// SizeOfNVarInt return the buffer size after encoding value as NVarInt
func sizeOfNVarInt(value int64, width int) int {
	const unit = 8 // bit width of encoding unit

	var lead = value >> (width - 1)
	for size := width / unit - 1; size > 0; size-- {
		var lookAhead = value >> (size*unit - 1)
		if lookAhead != lead {
			return size + 1
		}
	}
	return 1
}

func (codec *VarCodec) encodeNVarInt(buffer []byte, value int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}

	const unit = 8 // bit width of encoding unit
	for codec.Size > 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}

		codec.Size--
		buffer[codec.Ptr] = byte(value >> (codec.Size * unit))
		codec.Ptr++
	}
	return nil
}

func (codec *VarCodec) decodeNVarInt(buffer []byte, value *int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to decode")
	}

	const unit = 8      // bit width of encoding unit
	if codec.Size > 0 { // initialize sign bit
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
		*value = int64(int8(buffer[codec.Ptr]) >> (unit - 1))
		codec.Size = -codec.Size
	}

	for codec.Size < 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}

		codec.Size++
		*value = (*value << unit) | int64(buffer[codec.Ptr])
		codec.Ptr++
	}
	return nil
}
