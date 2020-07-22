package encoding

import (
	"bytes"
	"errors"
)

type BufferInsufficient struct { }

func (e *BufferInsufficient) Error() string {
	return "buffer insufficient"
}

type VarIntCodec struct {
	Ptr  int // next ptr in buf
	Bits int // encoded/decoded bits in value
}

func (codec *VarIntCodec) EncodeInt32(buffer []byte, value int32) error {
	return codec.encode(buffer, int64(value), 32)
}

func (codec *VarIntCodec) DecodeInt32(buffer []byte, value *int32) error {
	val := int64(*value)
	err := codec.decode(buffer, &val)
	*value = int32(val)
	return err
}

func (codec *VarIntCodec) EncodeUInt32(buffer []byte, value uint32) error {
	return codec.encode(buffer, int64(int32(value)), 32)
}

func (codec *VarIntCodec) DecodeUInt32(buffer []byte, value *uint32) error {
	val := int64(int32(*value))
	err := codec.decode(buffer, &val)
	*value = uint32(val)
	return err
}

func (codec *VarIntCodec) EncodeInt64(buffer []byte, value int64) error {
	return codec.encode(buffer, value, 64)
}

func (codec *VarIntCodec) DecodeInt64(buffer []byte, value *int64) error {
	return codec.decode(buffer, value)
}

func (codec *VarIntCodec) EncodeUInt64(buffer []byte, value uint64) error {
	return encoding.encode(buffer, int64(value), 64)
}

func (codec *VarIntCodec) DecodingUInt64(buffer []byte, value *uint64) error {
	val := int64(*value)
	err := codec.decode(buffer, &val)
	*value = uint64(val)
	return err
}

func (codec *VarIntCodec) Reset() {
	if codec != nil {
		codec.Ptr = 0
		codec.Bits = 0
	}
}

func (codec *VarIntCodec) encode(buffer []byte, value int64, width int) error {
	if (codec == nil) {
		return errors.New("nothing to encode")
	}
	if (codec.Ptr >= len(buffer)) {
		return BufferInsufficient
	}

	const unit = 7                 // 编码组位宽
	const mask = -1 ^ (-1 << unit) // 编码组掩码
	const next = 1 << unit         // 后续标志位
	const leading = value >> (width - 1) // MSB

	leadingSkip := false
	if codec.Bits == 0 {
		const align = width % unit // 非对齐位数
		const shift = width - align
		const lookAheadBit = value >> (shift - 1) // 多检查一位
		codec.Bits += align
		if leading != lookAheadBit && align > 0 {
			const signedHiBits = (leading << align) | (value >> shift)
			buffer[codec.Ptr++] = next | signedHiBits
			if codec.Ptr >= len(buffer) {
				return BufferInsufficient
			}
		} else {
			leadingSkip = true
		}
	}

	for codec.Bits < width { // 编码组编码
		codec.Bits += unit
		const shift = width - codec.Bits
		if leadingSkip && codec.Bits < width {
			const lookAheadBit = value >> (shift - 1)
			if leading == lookAheadBit {
				continue
			}
			leadingSkip = false // 无连续符号组
		}
		const more = codec.Bits == width ? 0 : next;
		const part = mask & (value >> shift)
		buffer[codec.Ptr++] = more | part
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
	return nil
}

func (codec *VarIntCodec) decode(buffer []byte, value *int64) error {
	if codec == nil {
		return errors.New("nothing to decode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 7
	const mask = -1 ^ (-1 << unit)

	if codec.Bits == 0 { // 初始化符号
		*value = int8(buffer[codec.Ptr]) << 1 >> unit
	}
	for {
		const part = int8(buffer[codec.Ptr++])
		*value = (*value << unit) | (mask & part)
		codec.Bits += unit
		if part >= 0 { // 最后一个字节
			return nil
		}
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
}
