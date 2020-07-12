package y3

import (
	"testing"
)

// 每个Packet最小长度是4个bytes
func TestLackLengthPacket(t *testing.T) {
	buf := []byte{0x01, 0x01, 0x01}
	expected := "invalid y3 packet minimal size"
	_, _, err := DecodePrimitivePacket(buf)
	if err.Error() != expected {
		t.Errorf("err actual = %v, and Expected = %v", err, expected)
	}
}

//
func TestPacketWrongLength(t *testing.T) {
	buf := []byte{0x04, 0x03, 0x02, 0x01}
	expected := "malformed, Length can not smaller than 2"
	_, _, err := DecodePrimitivePacket(buf)
	if err != nil && err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err)
	}
}

// TestUnknownType 测试还未定义的基础数据类型Type
func TestUnknownType(t *testing.T) {
	// 0x08 是未定义的Type
	buf := []byte{0x04, 0x04, 0x08, 0x01}
	expected := "Invalid PrimitiveType"

	_, _, err := DecodePrimitivePacket(buf)
	if err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err.Error())
	}
}

// 测试读取 0x04:-1
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x04, 0x01, 0x01}
	expectedTag := byte(0x04)
	var expectedLength int64 = 1
	expectedType := PrimitiveType(Varint)
	expectedValue := []byte{0x01}

	res, endPos, err := DecodePrimitivePacket(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res.Tag != expectedTag {
		t.Errorf("res.Tag actual = %v, and Expected = %v", res.Tag, expectedTag)
	}

	if res.Length() != expectedLength {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.Length(), expectedLength)
	}

	if res.Type != expectedType {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.Type, expectedType)
	}

	if !_compareByteSlice(res.basePacket.raw, expectedValue) {
		t.Errorf("res.raw actual = %v, and Expected = %v", res.basePacket.raw, expectedType)
	}

	if endPos != 4 {
		t.Errorf("endPos actual = %v, and Expected = %v", endPos, 4)
	}
}

// 测试读取 0x0A:-1
func TestParseInt64(t *testing.T) {
	buf := []byte{0x0A, 0x04, 0x01, 0x01}
	expected := int64(-1)

	res, _, err := DecodePrimitivePacket(buf)
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

// 测试 0x0B:"C"
func TestParseString(t *testing.T) {
	buf := []byte{0x0B, 0x04, 0x00, 0x43}
	expectedType := PrimitiveType(String)
	expectedValue := "C"

	res, _, err := DecodePrimitivePacket(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if expectedType != res.Type {
		t.Errorf("res.Type actual = %4b, and Expected = %4b", res.Type, expectedType)
	}

	target, err := res.ToUTF8String()
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if expectedValue != target {
		t.Errorf("Result actual = %v, and Expected = %v", t, expectedValue)
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
