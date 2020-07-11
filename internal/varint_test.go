package varint

import (
	"testing"
)

func TestOneByte(t *testing.T) {
	buf := []byte{0x01}
	dec, len := NewDecoder(buf)
	val, err := dec.Decode()
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if len != 1 {
		t.Errorf("l should %v, actual = %v", 1, len)
	}

	if val != -1 {
		t.Errorf("value should %v, actual = %v", -1, val)
	}
}

func Test2Byte(t *testing.T) {
	buf := []byte{0x80, 0x01}
	var expected int64 = 64
	expectedLength := 2
	dec, len := NewDecoder(buf)
	val, err := dec.Decode()
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if len != expectedLength {
		t.Errorf("l should %v, actual = %v", expectedLength, len)
	}

	if val != expected {
		t.Errorf("value should %v, actual = %v", -1, val)
	}
}

func TestZigZagEncode(t *testing.T) {
	data := []int64{0, -1, 1, -2, 2, -3, 3}
	expected := []uint64{0, 1, 2, 3, 4, 5, 6}
	// if zigzagDecode(1) != -1 {}
	for i, v := range data {
		tmp := zigzagEncode(v)
		if tmp != expected[i] {
			t.Errorf("zigzad encode %v should %v, actual is %v", v, expected[i], tmp)
		}
	}
}

func TestZigZagDecode(t *testing.T) {
	data := []uint64{0, 1, 2, 3, 4, 5, 6}
	expected := []int64{0, -1, 1, -2, 2, -3, 3}
	// if zigzagDecode(1) != -1 {}
	for i, v := range data {
		tmp := zigzagDecode(v)
		if tmp != expected[i] {
			t.Errorf("zigzad decode %v should %v, actual is %v", v, expected[i], tmp)
		}
	}
}
