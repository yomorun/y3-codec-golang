package y3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 每个Packet最小长度是2个bytes
func TestLackLengthPacket(t *testing.T) {
	buf := []byte{0x01}
	expected := "invalid y3 packet minimal size"
	_, _, _, err := DecodePrimitivePacket(buf)
	if err.Error() != expected {
		t.Errorf("err actual = %v, and Expected = %v", err, expected)
	}
}

func TestPacketWrongLength(t *testing.T) {
	buf := []byte{0x04, 0x00, 0x02, 0x01}
	expected := "malformed, Length can not smaller than 1"
	_, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	if err != nil && err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err)
	}
}

// 测试读取 0x04:-1
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x01, 0x7F}
	expectedTag := byte(0x04)
	var expectedLength uint32 = 1
	expectedValue := []byte{0x7F}

	res, endPos, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	assert.Equal(t, expectedTag, res.SeqID())
	assert.Equal(t, expectedLength, res.length)

	if !_compareByteSlice(res.valBuf, expectedValue) {
		t.Errorf("res.raw actual = %v, and Expected = %v", res.valBuf, expectedValue)
	}

	assert.Equal(t, endPos, 3)
}

// 测试读取 0x0A:2
func TestParseInt32(t *testing.T) {
	buf := []byte{0x0A, 0x02, 0x81, 0x7F}
	expectedValue := int32(255)

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// 测试 0x0B:"C"
func TestParseString(t *testing.T) {
	buf := []byte{0x0B, 0x01, 0x43}
	expectedValue := "C"

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToUTF8String()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// 测试 0x0B:"C"
func TestParseEmptyString(t *testing.T) {
	buf := []byte{0x0B, 0x00}
	expectedValue := ""

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToUTF8String()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// compares two slice, every element is equal
func _compareByteSlice(left []byte, right []byte) bool {
	if len(left) != len(right) {
		return false
	}

	for i, v := range left {
		if v != right[i] {
			return false
		}
	}

	return true
}
