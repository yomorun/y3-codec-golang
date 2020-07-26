package encoding

import (
	"errors"
)

// SizeOfPVarInt32 return the buffer size after encoding value as PVarInt32
func SizeOfPVarInt32(value int32) int {
	return sizeOfPVarInt(int64(value), 32)
}

// EncodePVarInt32 encode value as PVarInt32 to buffer
func (codec *VarCodec) EncodePVarInt32(buffer []byte, value int32) error {
	return codec.encodePVarInt(buffer, int64(value))
}

// DecodePVarInt32 decode to value as PVarInt32 from buffer
func (codec *VarCodec) DecodePVarInt32(buffer []byte, value *int32) error {
	var val = int64(*value)
	var err = codec.decodePVarInt(buffer, &val)
	*value = int32(val)
	return err
}

// SizeOfPVarUInt32 return the buffer size after encoding value as PVarUInt32
func SizeOfPVarUInt32(value uint32) int {
	return sizeOfPVarInt(int64(int32(value)), 32)
}

// EncodePVarUInt32 encode value as PVarUInt32 to buffer
func (codec *VarCodec) EncodePVarUInt32(buffer []byte, value uint32) error {
	return codec.encodePVarInt(buffer, int64(int32(value)))
}

// DecodePVarUInt32 decode to value as PVarUInt32 from buffer
func (codec *VarCodec) DecodePVarUInt32(buffer []byte, value *uint32) error {
	var val = int64(int32(*value))
	var err = codec.decodePVarInt(buffer, &val)
	*value = uint32(val)
	return err
}

// SizeOfPVarInt64 return the buffer size after encoding value as PVarInt64
func SizeOfPVarInt64(value int64) int {
	return sizeOfPVarInt(value, 64)
}

// EncodePVarInt64 encode value as PVarInt64 to buffer
func (codec *VarCodec) EncodePVarInt64(buffer []byte, value int64) error {
	return codec.encodePVarInt(buffer, value)
}

// DecodePVarInt64 decode to value as PVarInt64 from buffer
func (codec *VarCodec) DecodePVarInt64(buffer []byte, value *int64) error {
	return codec.decodePVarInt(buffer, value)
}

// SizeOfPVarUInt64 return the buffer size after encoding value as PVarUInt64
func SizeOfPVarUInt64(value uint64) int {
	return sizeOfPVarInt(int64(value), 64)
}

// EncodePVarUInt64 encode value as PVarUInt64 to buffer
func (codec *VarCodec) EncodePVarUInt64(buffer []byte, value uint64) error {
	return codec.encodePVarInt(buffer, int64(value))
}

// DecodePVarUInt64 decode to value as PVarUInt64 from buffer
func (codec *VarCodec) DecodePVarUInt64(buffer []byte, value *uint64) error {
	var val = int64(*value)
	var err = codec.decodePVarInt(buffer, &val)
	*value = uint64(val)
	return err
}

func sizeOfPVarInt(value int64, width int) int {
	const unit = 7 // bit width of encoding unit

	var lead = value >> (width - 1)
	for size := width / unit; size > 0; size-- {
		var lookAhead = value >> (size*unit - 1)
		if lookAhead != lead {
			return size + 1
		}
	}
	return 1
}

func (codec *VarCodec) encodePVarInt(buffer []byte, value int64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}
	if codec.Ptr >= len(buffer) {
		return ErrBufferInsufficient
	}

	const unit = 7       // bit width of encoding unit
	const more = -1 << 7 // continuation bits
	for codec.Size > 1 {
		codec.Size--
		var part = value >> (codec.Size * unit)

		buffer[codec.Ptr] = byte(part | more)
		codec.Ptr++

		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
	}

	const mask = -1 ^ (-1 << unit) // mask for encoding unit
	codec.Size = 0
	buffer[codec.Ptr] = byte(value & mask)
	codec.Ptr++
	return nil
}

func (codec *VarCodec) decodePVarInt(buffer []byte, value *int64) error {
	if codec == nil {
		return errors.New("nothing to decode")
	}
	if codec.Ptr >= len(buffer) {
		return ErrBufferInsufficient
	}

	const unit = 7 // bit width of encoding unit
	if codec.Size == 0 { // initialize sign bit
		const flag = 8 - unit // bit width for non-encoding bits
		*value = int64(int8(buffer[codec.Ptr]) << flag >> unit)
	}

	const mask = -1 ^ (-1 << unit) // mask for encoding unit
	for {
		var part = int8(buffer[codec.Ptr])
		codec.Ptr++

		codec.Size++
		*value = (*value << unit) | int64(mask & part)

		if part >= 0 { // it's the last byte
			return nil
		}
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
	}
}
