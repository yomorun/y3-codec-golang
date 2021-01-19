package y3

import (
	"testing"
)

// 每个Packet最小长度是3个bytes
func TestLackLengthPacket(t *testing.T) {
	buf := []byte{0x01, 0x01}
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
	if err != nil && err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err)
	}
}

// 测试读取 0x04:-1
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x01, 0x7F}
	expectedTag := byte(0x04)
	var expectedLength int32 = 1
	expectedValue := []byte{0x7F}

	res, endPos, _, err := DecodePrimitivePacket(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res.SeqID() != expectedTag {
		t.Errorf("res.Tag actual = %v, and Expected = %v", res.SeqID(), expectedTag)
	}

	if res.length != uint32(expectedLength) {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.length, expectedLength)
	}

	if !_compareByteSlice(res.valBuf, expectedValue) {
		t.Errorf("res.raw actual = %v, and Expected = %v", res.valBuf, expectedValue)
	}

	if endPos != 3 {
		t.Errorf("endPos actual = %v, and Expected = %v", endPos, 3)
	}
}

// 测试读取 0x0A:2
func TestParseInt32(t *testing.T) {
	// 原例子数据存在问题，修改如下：
	//buf := []byte{0x0A, 0x02, 0x01, 0x02}
	buf := []byte{0x0A, 0x02, 0x81, 0x7F}
	expected := int32(255)

	res, _, _, err := DecodePrimitivePacket(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	val, err := res.ToInt32()
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if val != expected {
		t.Errorf("value actual = %v, and Expected = %v", val, expected)
	}
}

// 测试 0x0B:"C"
func TestParseString(t *testing.T) {
	buf := []byte{0x0B, 0x01, 0x43}
	expectedValue := "C"

	res, _, _, err := DecodePrimitivePacket(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
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
