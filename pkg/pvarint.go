package encoding

import (
	"errors"
)

var BufferInsufficient = errors.New("buffer insufficient")

type PVarIntCodec struct {
	Ptr  int // next ptr in buf
	Bits int // encoded/decoded bits in value
}

func (codec *PVarIntCodec) EncodeInt32(buffer []byte, value int32) error {
	return codec.encode(buffer, int64(value), 32)
}

func (codec *PVarIntCodec) DecodeInt32(buffer []byte, value *int32) error {
	var val = int64(*value)
	var err = codec.decode(buffer, &val)
	*value = int32(val)
	return err
}

func (codec *PVarIntCodec) EncodeUInt32(buffer []byte, value uint32) error {
	return codec.encode(buffer, int64(int32(value)), 32)
}

func (codec *PVarIntCodec) DecodeUInt32(buffer []byte, value *uint32) error {
	var val = int64(int32(*value))
	var err = codec.decode(buffer, &val)
	*value = uint32(val)
	return err
}

func (codec *PVarIntCodec) EncodeInt64(buffer []byte, value int64) error {
	return codec.encode(buffer, value, 64)
}

func (codec *PVarIntCodec) DecodeInt64(buffer []byte, value *int64) error {
	return codec.decode(buffer, value)
}

func (codec *PVarIntCodec) EncodeUInt64(buffer []byte, value uint64) error {
	return codec.encode(buffer, int64(value), 64)
}

func (codec *PVarIntCodec) DecodeUInt64(buffer []byte, value *uint64) error {
	var val = int64(*value)
	var err = codec.decode(buffer, &val)
	*value = uint64(val)
	return err
}

func (codec *PVarIntCodec) Reset() {
	if codec != nil {
		codec.Ptr = 0
		codec.Bits = 0
	}
}

func (codec *PVarIntCodec) encode(buffer []byte, value int64, width int) error {
	if codec == nil {
		return errors.New("nothing to encode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 7          // 编码组位宽
	const next = -1 << unit // 接续标志位
	const mask = -1 ^ next  // 编码组掩码
	const leading = value >> (width - 1) // 符号位符号扩展

	var leadingSkip = false
	if codec.Bits == 0 {
		var align = width % unit  // 非对齐位数
		var shift = width - align
		var lookAheadBit = value >> (shift - 1) // 多检查一位

		codec.Bits += align
		if leading != lookAheadBit {
			var signedHiBits = (leading << align) | (value >> shift)
			buffer[codec.Ptr] = byte(next | signedHiBits)

			codec.Ptr++
			if codec.Ptr >= len(buffer) {
				return BufferInsufficient
			}
		} else {
			leadingSkip = true
		}
	}

	for codec.Bits < width { // 编码组编码
		codec.Bits += unit
		var shift = width - codec.Bits
		if leadingSkip && codec.Bits < width {
			var lookAheadBit = value >> (shift - 1)
			if leading == lookAheadBit {
				continue
			}
			leadingSkip = false // 无连续符号组
		}

		var part = mask & (value >> shift)
		if codec.Bits == width {
			buffer[codec.Ptr] = byte(part)
			codec.Ptr++
			return nil
		}

		buffer[codec.Ptr] = byte(next | part)
		codec.Ptr++
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
}

func (codec *PVarIntCodec) decode(buffer []byte, value *int64) error {
	if codec == nil {
		return errors.New("nothing to decode")
	}
	if codec.Ptr >= len(buffer) {
		return BufferInsufficient
	}

	const unit = 7
	const mask = -1 ^ (-1 << unit)

	if codec.Bits == 0 { // 初始化符号
		*value = int64(int8(buffer[codec.Ptr]) << 1 >> 7)
	}

	for {
		var part = int8(buffer[codec.Ptr])
		*value = (*value << unit) | int64(mask & part)

		codec.Ptr++
		codec.Bits += unit
		if part >= 0 { // 最后一个字节
			return nil
		}
		if codec.Ptr >= len(buffer) {
			return BufferInsufficient
		}
	}
}

func EncodePvarint(val int32) (buf []byte, length int, err error) {
	var c = PVarIntCodec{}
	buf = make([]byte, 10)
	err = c.EncodeInt32(buf, val)
	return buf, len(buf), err
}

func Pvarint(buf []byte, pos int) (res int32, length int, err error) {
	var c = PVarIntCodec{}
	var r *int32
	err = c.DecodeInt32(buf, r)
	return *r, 0, err
}

func EncodeUpvarint(val uint32) (buf []byte, length int, err error) {
	var c = PVarIntCodec{}
	buf = make([]byte, 10)
	err = c.EncodeUInt32(buf, val)
	return buf, len(buf), err
}

func Upvarint(buf []byte, pos int) (res uint32, length int, err error) {
	var c = PVarIntCodec{}
	var r *uint32
	err = c.DecodeUInt32(buf, r)
	return uint32(*r), 0, err
}
