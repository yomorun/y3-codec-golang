package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPInt32(t *testing.T) {
	testPVarInt32(t, -1, []byte{0x7F})
	testPVarInt32(t, -5, []byte{0x7B})
	testPVarInt32(t, 63, []byte{0x3F})
	testPVarInt32(t, -65, []byte{0xFF, 0x3F})
	testPVarInt32(t, 127, []byte{0x80, 0x7F})
	testPVarInt32(t, 255, []byte{0x81, 0x7F})
	testPVarInt32(t, -4097, []byte{0xDF, 0x7F})
	testPVarInt32(t, -8193, []byte{0xFF, 0xBF, 0x7F})
	testPVarInt32(t, -2097152, []byte{0xFF, 0x80, 0x80, 0x00})
	testPVarInt32(t, -134217729, []byte{0xFF, 0xBF, 0xFF, 0xFF, 0x7F})
	testPVarInt32(t, -2147483648, []byte{0xF8, 0x80, 0x80, 0x80, 0x00})
}

func TestPUInt32(t *testing.T) {
	testPVarUInt32(t, 1, []byte{0x01})
	testPVarUInt32(t, 127, []byte{0x80, 0x7F})
	testPVarUInt32(t, 128, []byte{0x81, 0x00})
	testPVarUInt32(t, 130, []byte{0x81, 0x02})
	testPVarUInt32(t, 1048576, []byte{0x80, 0xC0, 0x80, 0x00})
	testPVarUInt32(t, 134217728, []byte{0x80, 0xC0, 0x80, 0x80, 0x00})
	testPVarUInt32(t, 4294967295, []byte{0x7F})
}

func TestPInt64(t *testing.T) {
	testPVarInt64(t, 0, []byte{0x00})
	testPVarInt64(t, 1, []byte{0x01})
	testPVarInt64(t, -1, []byte{0x7F})
}

func TestPUInt64(t *testing.T) {
	testPVarUInt64(t, 0, []byte{0x00})
	testPVarUInt64(t, 1, []byte{0x01})
	testPVarUInt64(t, 18446744073709551615, []byte{0x7F})
}

func testPVarInt32(t *testing.T, value int32, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, uint32(value), bytes)
	var size = SizeOfPVarInt32(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodePVarInt32(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val int32
	codec = VarCodec{}
	assert.Nil(t, codec.DecodePVarInt32(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testPVarUInt32(t *testing.T, value uint32, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, value, bytes)
	var size = SizeOfPVarUInt32(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodePVarUInt32(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val uint32
	codec = VarCodec{}
	assert.Nil(t, codec.DecodePVarUInt32(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testPVarInt64(t *testing.T, value int64, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, uint64(value), bytes)
	var size = SizeOfPVarInt64(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodePVarInt64(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val int64
	codec = VarCodec{}
	assert.Nil(t, codec.DecodePVarInt64(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testPVarUInt64(t *testing.T, value uint64, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, value, bytes)
	var size = SizeOfPVarUInt64(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodePVarUInt64(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val uint64
	codec = VarCodec{}
	assert.Nil(t, codec.DecodePVarUInt64(bytes, &val), msg)
	assert.Equal(t, value, val)
}
