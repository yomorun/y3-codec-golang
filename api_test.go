package y3

import (
	"testing"
)

func TestD2EncodeUInt32(t *testing.T) {
	testD2EncodeUInt32(t, 0x03, 6, []byte{0x03, 0x01, 0x06})
	testD2EncodeUInt32(t, 0x06, 127, []byte{0x06, 0x02, 0x80, 0x7F})
}

func TestD2EncodeInt32(t *testing.T) {
	testD2EncodeInt(t, 0x03, -1, []byte{0x03, 0x01, 0x7F})
	testD2EncodeInt(t, 0x06, -65, []byte{0x06, 0x02, 0xFF, 0x3F})
	testD2EncodeInt(t, 0x09, 255, []byte{0x09, 0x02, 0x81, 0x7F})
}

func TestD2EncodeUInt64(t *testing.T) {
	testD2EncodeUInt64(t, 0x03, 0, []byte{0x03, 0x01, 0x00})
	testD2EncodeUInt64(t, 0x06, 1, []byte{0x06, 0x01, 0x01})
	testD2EncodeUInt64(t, 0x09, 18446744073709551615, []byte{0x09, 0x01, 0x7F})
}

func TestD2EncodeInt64(t *testing.T) {
	testD2EncodeInt64(t, 0x03, 0, []byte{0x03, 0x01, 0x00})
	testD2EncodeInt64(t, 0x06, 1, []byte{0x06, 0x01, 0x01})
	testD2EncodeInt64(t, 0x09, -1, []byte{0x09, 0x01, 0x7F})
}

func TestD2EncodeFloat32(t *testing.T) {
	testD2EncodeFloat32(t, 0x03, -2, []byte{0x03, 0x01, 0xC0})
	testD2EncodeFloat32(t, 0x06, 0.25, []byte{0x06, 0x02, 0x3E, 0x80})
	testD2EncodeFloat32(t, 0x09, 68.123, []byte{0x09, 0x04, 0x42, 0x88, 0x3E, 0xFA})
}

func TestD2EncodeFloat64(t *testing.T) {
	testD2EncodeFloat64(t, 0x03, 23, []byte{0x03, 0x02, 0x40, 0x37})
	testD2EncodeFloat64(t, 0x06, 2, []byte{0x06, 0x01, 0x40})
	testD2EncodeFloat64(t, 0x09, 0.01171875, []byte{0x09, 0x02, 0x3F, 0x88})
}

func TestD2EncodeString(t *testing.T) {
	p, _ := EncodeString(0x01, "C")
	compareTwoBytes(t, p, []byte{0x01, 0x01, 0x43})
	p, _ = EncodeString(0x01, "CC")
	compareTwoBytes(t, p, []byte{0x01, 0x02, 0x43, 0x43})
	p, _ = EncodeString(0x01, "Yona")
	compareTwoBytes(t, p, []byte{0x01, 0x04, 0x59, 0x6F, 0x6E, 0x61})
	p, _ = EncodeString(0x01, "https://yomo.run")
	compareTwoBytes(t, p, []byte{0x01, 0x10, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3A, 0x2F, 0x2F, 0x79, 0x6F, 0x6D, 0x6F, 0x2E, 0x72, 0x75, 0x6E})
}

func TestD2EncodeBytes(t *testing.T) {
	p, _ := EncodeBytes(0x01, []byte{0x03, 0x06, 0x09, 0x0C, 0x0F})
	compareTwoBytes(t, p, []byte{0x01, 0x05, 0x03, 0x06, 0x09, 0x0C, 0x0F})
}

func TestD2EncodeBool(t *testing.T) {
	p, _ := EncodeBool(0x01, true)
	compareTwoBytes(t, p, []byte{0x01, 0x01, 0x01})
	p, _ = EncodeBool(0x01, false)
	compareTwoBytes(t, p, []byte{0x01, 0x01, 0x00})
}

func testD2EncodeUInt32(t *testing.T, tag int, val uint, expected []byte) {
	p, _ := EncodeUInt(tag, val)
	compareTwoBytes(t, p, expected)
}

func testD2EncodeInt(t *testing.T, tag int, val int, expected []byte) {
	p, _ := EncodeInt(tag, val)
	compareTwoBytes(t, p, expected)
}

func testD2EncodeUInt64(t *testing.T, tag int, val uint64, expected []byte) {
	p, _ := EncodeUInt64(tag, val)
	compareTwoBytes(t, p, expected)
}

func testD2EncodeInt64(t *testing.T, tag int, val int64, expected []byte) {
	p, _ := EncodeInt64(tag, val)
	compareTwoBytes(t, p, expected)
}

func testD2EncodeFloat32(t *testing.T, tag int, val float32, expected []byte) {
	p, _ := EncodeFloat32(tag, val)
	compareTwoBytes(t, p, expected)
}

func testD2EncodeFloat64(t *testing.T, tag int, val float64, expected []byte) {
	p, _ := EncodeFloat64(tag, val)
	compareTwoBytes(t, p, expected)
}

func compareTwoBytes(t *testing.T, result []byte, expected []byte) {
	for i, p := range result {
		if p != expected[i] {
			t.Errorf("\nexpected:[% X]\n  actual:[% X]\n", expected, result)
			break
		}
	}
}
