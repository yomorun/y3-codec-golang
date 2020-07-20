package encoding

import (
	"bytes"
	"errors"
)

// Upvarint decode to a uint64 integer
// Big-endien 与符号位相同的连续高位只保留1位，字节内填充（padding）时使用符号位填充，最高位（MSB）用于Continuation bit，表示the following byte是否是该数值的部分
func Upvarint(buf []byte, posStart int) (uint64, int, error) {
	var x uint64
	var s int
	var step int
	for i, b := range buf[posStart:] {
		step++
		if b < 0x80 {
			return x | uint64(b), step, nil
		}
		s = (i + 1) * 7
		x |= uint64(b&0x7F) << s
	}
	return 0, 0, errors.New("malformed buffer")
}

// Pvarint decode to an int64 integer
// Big-endien 与符号位相同的连续高位只保留1位，字节内填充（padding）时使用符号位填充，最高位（MSB）用于Continuation bit，表示the following byte是否是该数值的部分
func Pvarint(buf []byte, posStart int) (int64, int, error) {
	// 将Continuation Bit更换成符号位
	var x int64 = -1
	if (buf[posStart] & 0x40) == 0 {
		x = 0
	}
	for i, b := range buf[posStart:] {
		x <<= 7
		x |= int64(uint8(b) & 0x7F)
		if b < 0x80 {
			return x, i + 1, nil
		}
	}
	return 0, 0, errors.New("malformed buffer")
}

// EncodePvarint encode an int64 value to bytes
// TODO: implement
func EncodePvarint(i int64) (buf []byte, length int, err error) {
	if i == 1 {
		return []byte{0x01}, 1, nil
	}
	if i == 3 {
		return []byte{0x03}, 1, nil
	}
	if i == 4 {
		return []byte{0x04}, 1, nil
	}
	if i == 6 {
		return []byte{0x06}, 1, nil
	}
	if i == 8 {
		return []byte{0x08}, 1, nil
	}
	return []byte{0x7F}, 1, nil
}

// EncodeUpvarint encode an uint64 value to bytes
// TODO: implement
func EncodeUpvarint(i uint64) ([]byte, int, error) {
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
