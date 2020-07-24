package encoding

import (
	"bytes"
	"errors"
)

func SizeOfInt32(value int32) int {
	return sizeOf(int64(value), 32)
}

func EncodeInt32(buffer []byte, size int, value int32) {
	encode(buffer, size, int64(value))
}

func DecodeInt32(buffer []byte, size int) int32 {
	return int32(decode(buffer, size))
}

func SizeOfUInt32(value uint32) int {
	return sizeOf(int64(int32(value)), 32)
}

func EncodeUInt32(buffer []byte, size int, value uint32) {
	encode(buffer, size, int64(int32(value)))
}

func DecodeUInt32(buffer []byte, size int) uint32 {
	return uint32(decode(buffer, size))
}

func SizeOfInt64(value int64) int {
	return sizeOf(value, 64)
}

func EncodeInt64(buffer []byte, size int, value int64) {
	encode(buffer, size, value)
}

func DecodeInt64(buffer []byte, size int) int64 {
	return decode(buffer, size)
}

func SizeOfUInt64(value uint64) int {
	return sizeOf(int64(value), 64)
}

func EncodeUInt64(buffer []byte, size int, value uint64) {
	encode(buffer, size, int64(value))
}

func DecodingUInt64(buffer []byte, size int) uint64 {
	return uint64(decode(buffer, size))
}

func sizeOf(value int64, width int) int {
	const unit = 8
	var size = width / unit

	var shift = width - unit
	const lead = value >> shift

	for shift > 0 {
		shift -= unit
		if (lead != (value >> shift)) {
			break
		}
		--size
	}
	return size
}

func encode(buffer []byte, size int, value int64) {
	const unit = 8
	for i := 0; size > 0; ++i {
		--size
		buffer[i] = byte(value >> (size * unit))
	}
}

func decode(buffer []byte, size int) int64 {
	const unit = 8
	var value = int64(int8(buffer[0]) >> 7)
	for i := 0; i < size; ++i {
		value = (value << unit) | int64(buffer[i])
	}
	return value
}
