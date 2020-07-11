package y3

import (
	"testing"
)

// 每个Packet最小长度是4个bytes
func TestLackLengthPacket(t *testing.T) {
	buf := []byte{0x01, 0x01, 0x01}
	expected := "invalid y3 packet minimal size"
	_, err := Decode(buf)
	if err.Error() != expected {
		t.Errorf("err actual = %v, and Expected = %v", err, expected)
	}
}

//
func TestPacketWrongLength(t *testing.T) {
	buf := []byte{0x04, 0x03, 0x02, 0x01}
	expected := "malformed, Length can not smaller than 2"
	_, err := Decode(buf)
	if err != nil && err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err)
	}
}

// TestUnknownType 测试还未定义的基础数据类型Type
func TestUnknownType(t *testing.T) {
	// 0x08 是未定义的Type
	buf := []byte{0x04, 0x04, 0x08, 0x01}
	expected := "Invalid Type"

	_, err := Decode(buf)
	if err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err.Error())
	}
}

// 测试读取 0x04:-1
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x04, 0x01, 0x01}
	expectedTag := byte(0x04)
	var expectedLength int64 = 1
	expectedType := Type(Varint)
	expectedValue := []byte{0x01}

	res, err := Decode(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res.Tag != expectedTag {
		t.Errorf("res.Tag actual = %v, and Expected = %v", res.Tag, expectedTag)
	}

	if res.Length != expectedLength {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.Length, expectedLength)
	}

	if res.Type != expectedType {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.Type, expectedType)
	}

	if !_compareByteSlice(res.raw, expectedValue) {
		t.Errorf("res.raw actual = %v, and Expected = %v", res.raw, expectedType)
	}
}

// 测试读取 0x0A:-1
func TestInt64PacketRead(t *testing.T) {
	buf := []byte{0x0A, 0x04, 0x01, 0x01}
	expected := int64(-1)

	res, err := Decode(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	val, err := res.ToInt64()
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if val != expected {
		t.Errorf("value actual = %v, and Expected = %v", val, expected)
	}
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
