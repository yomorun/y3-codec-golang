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

// TestUnknownType 测试还未定义的基础数据类型Type
func TestUnknownType(t *testing.T) {
	// 0x08 是未定义的Type
	buf := []byte{0x04, 0x02, 0x08, 0x01}
	expected := "Invalid Type"

	_, err := Decode(buf)
	if err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err.Error())
	}
}

// 读取
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x02, 0x01, 0x01}
	expectedTag := byte(0x04)
	expectedLength := 2
	expectedType := Type(Varint)

	res, err := Decode(buf)
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res.Tag != expectedTag {
		t.Errorf("res.Tag actual = %v, and Expected = %v", res.Tag, expectedTag)
	}

	if res.Length != 2 {
		t.Errorf("res.Length actual = %v, and Expected = %v", res, expectedLength)
	}

	if res.Type != expectedType {
		t.Errorf("res.Length actual = %v, and Expected = %v", res.Type, expectedType)
	}
}
