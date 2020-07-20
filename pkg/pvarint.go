package encoding

import (
	"bytes"
	"errors"
)

// Upvarint decode to a uint32 integer
// Big-endien 与符号位相同的连续高位只保留1位，字节内填充（padding）时使用符号位填充，最高位（MSB）用于Continuation bit，表示the following byte是否是该数值的部分
func Upvarint(buf []byte, posStart int) (uint32, int, error) {
	var x uint32
	var s int
	var step int
	for i, b := range buf[posStart:] {
		step++
		if b < 0x80 {
			return x | uint32(b), step, nil
		}
		s = (i + 1) * 7
		x |= uint32(b&0x7F) << s
	}
	return 0, 0, errors.New("malformed buffer")
}

// Pvarint decode to an int32 integer
// Big-endien 与符号位相同的连续高位只保留1位，字节内填充（padding）时使用符号位填充，最高位（MSB）用于Continuation bit，表示the following byte是否是该数值的部分
func Pvarint(buf []byte, posStart int) (int32, int, error) {
	// 将Continuation Bit更换成符号位
	var x int32 = -1
	if (buf[posStart] & 0x40) == 0 {
		x = 0
	}
	for i, b := range buf[posStart:] {
		x <<= 7
		x |= int32(uint8(b) & 0x7F)
		if b < 0x80 {
			return x, i + 1, nil
		}
	}
	return 0, 0, errors.New("malformed buffer")
}

// EncodePvarint encode an int32 value to bytes
func EncodePvarint(i int32) ([]byte, int, error) {
	if i >= 0 {
		return EncodeUpvarint(uint32(i))
	}
	panic("!!!!!!!!!!!!!!!!!!!")
}

// EncodeUpvarint encode an uint32 value to bytes
func EncodeUpvarint(i uint32) ([]byte, int, error) {
	buf := new(bytes.Buffer)
	len := 0
	for {
		if i == 0 {
			if len == 0 {
				buf.WriteByte(0x00)
			}
			return reverse(buf.Bytes()), buf.Len(), nil
		}
		p := i & 0x7F
		if len == 0 {
			buf.WriteByte(byte(p))
		} else {
			buf.WriteByte(byte(p | 0x80))
		}
		i = i >> 7
		len++
	}
}

func reverse(b []byte) (r []byte) {
	l := len(b)
	r = make([]byte, l)
	for i, v := range b {
		r[l-i-1] = v
	}
	return r
}
