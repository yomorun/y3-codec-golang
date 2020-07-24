package encoding

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
	var lead = value >> (width - 1)

	var size = width / 8
	for s := size - 1; s > 0; s-- {
		if lead == (value >> (s * 8)) {
			size--
		}
	}
	return size
}

func encode(buffer []byte, size int, value int64) {
	for i := 0; size > 0; i++ {
		size--
		buffer[i] = byte(value >> (size * 8))
	}
}

func decode(buffer []byte, size int) int64 {
	var value = int64(int8(buffer[0]) >> 7)
	for i := 0; i < size; i++ {
		value = (value << 8) | int64(buffer[i])
	}
	return value
}
