package encoding

import (
	"errors"
	"math"
)

// SizeOfVarFloat32 return the buffer size after encoding value as VarFloat32
func SizeOfVarFloat32(value float32) int {
	return sizeOfVarFloat(uint64(math.Float32bits(value)), 32)
}

// EncodeVarFloat32 encode value as VarFloat32 to buffer
func (codec *VarCodec) EncodeVarFloat32(buffer []byte, value float32) error {
	return codec.encodeVarFloat(buffer, uint64(math.Float32bits(value)))
}

// DecodeVarFloat32 decode to value as VarFloat32 from buffer
func (codec *VarCodec) DecodeVarFloat32(buffer []byte, value *float32) error {
	var bits = uint64(math.Float32bits(*value))
	var err = codec.decodeVarFloat(buffer, &bits)
	*value = math.Float32frombits(uint32(bits))
	return err
}

// SizeOfVarFloat64 return the buffer size after encoding value as VarFloat32
func SizeOfVarFloat64(value float64) int {
	return sizeOfVarFloat(math.Float64bits(value), 64)
}

// EncodeVarFloat64 encode value as VarFloat64 to buffer
func (codec *VarCodec) EncodeVarFloat64(buffer []byte, value float64) error {
	return codec.encodeVarFloat(buffer, math.Float64bits(value))
}

// DecodeVarFloat64 decode to value as VarFloat64 from buffer
func (codec *VarCodec) DecodeVarFloat64(buffer []byte, value *float64) error {
	var bits = math.Float64bits(*value)
	var err = codec.decodeVarFloat(buffer, &bits)
	*value = math.Float64frombits(bits)
	return err
}

func sizeOfVarFloat(bits uint64, width int) int {
	const unit = 8    // bit width of encoding unit
	const mask = 0xFF // mask of encoding unit

	var size = width / unit
	for s := 0; size > 1; s += unit {
		if bits&(mask<<s) != 0 {
			return size
		}
		size--
	}
	return 1
}

func (codec *VarCodec) encodeVarFloat(buffer []byte, bits uint64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}

	const unit = 8 // bit width of encoding unit
	for codec.Size > 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
		codec.Size--
		buffer[codec.Ptr] = byte(bits >> (codec.Size * unit))
		codec.Ptr++
	}
	return nil
}

func (codec *VarCodec) decodeVarFloat(buffer []byte, bits *uint64) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to decode")
	}

	const unit = 8 // bit width of encoding unit
	for codec.Size > 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
		codec.Size--
		*bits = (*bits << unit) | uint64(buffer[codec.Ptr])
		codec.Ptr++
	}
	//TODO: fill 0 to the gap between encoded size and decoded size
	return nil
}
