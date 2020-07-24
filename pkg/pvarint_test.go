package encoding

import (
	"testing"
)

func TestU1Byte(t *testing.T) {
	buf := []byte{0x01}
	res, step, err := Upvarint(buf, 0)
	var expected uint32 = 1
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 1 {
		t.Errorf("step should %v, actual = %v", 1, step)
	}
}

func TestUWrongByte(t *testing.T) {
	buf := []byte{0x81, 0x82}
	_, _, err := Upvarint(buf, 0)
	if err == nil {
		t.Errorf("err should not nil")
	}
	if err.Error() != "malformed buffer" {
		t.Errorf("err should not nil, actual err = %v", err)
	}
}

func TestU2Bytes(t *testing.T) {
	buf := []byte{0x81, 0x02}
	res, step, err := Upvarint(buf, 0)
	var expected uint32 = 130
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 2 {
		t.Errorf("step should %v, actual = %v", 2, step)
	}
}

func TestU2BytesSkip(t *testing.T) {
	buf := []byte{0x81, 0x82, 0x81, 0x02}
	res, step, err := Upvarint(buf, 2)
	var expected uint32 = 130
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}

	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 2 {
		t.Errorf("step should %v, actual = %v", 2, step)
	}
}

func Test1ByteNagitive(t *testing.T) {
	buf := []byte{0x7B}
	res, step, err := Pvarint(buf, 0)
	var expected int32 = -5
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 1 {
		t.Errorf("step should %v, actual = %v", 1, step)
	}
}

func Test1BytePositive(t *testing.T) {
	buf := []byte{0x3F}
	res, step, err := Pvarint(buf, 0)
	var expected int32 = 63
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 1 {
		t.Errorf("step should %v, actual = %v", 1, step)
	}
}

func Test2BytesNagitive(t *testing.T) {
	buf := []byte{0xFF, 0x3F}
	res, step, err := Pvarint(buf, 0)
	var expected int32 = -65
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 2 {
		t.Errorf("step should %v, actual = %v", 2, step)
	}
}

func Test2BytesPositive2(t *testing.T) {
	buf := []byte{0x80, 0x7F}
	res, step, err := Pvarint(buf, 0)
	var expected int32 = 127
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 2 {
		t.Errorf("step should %v, actual = %v", 2, step)
	}
}

func Test2BytesPositive(t *testing.T) {
	buf := []byte{0xFF, 0x81, 0x7F}
	res, step, err := Pvarint(buf, 1)
	var expected int32 = 255
	if err != nil {
		t.Errorf("err should nil, actual = %v", err)
	}
	if res != expected {
		t.Errorf("value should %v, actual = %v", expected, res)
	}
	if step != 2 {
		t.Errorf("step should %v, actual = %v", 2, step)
	}
}

func TestUpvarintEncode1(t *testing.T) {
	var val uint32 = 127
	expected := []byte{0x7F}
	res, length, err := EncodeUpvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 1 {
		t.Errorf("expected length=1, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestUpvarintEncode2(t *testing.T) {
	var val uint32 = 128
	expected := []byte{0x81, 0x00}
	res, length, err := EncodeUpvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 2 {
		t.Errorf("expected length=2, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestUpvarintEncode3(t *testing.T) {
	var val uint32 = 1048576
	expected := []byte{0xC0, 0x80, 0x00}
	res, length, err := EncodeUpvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 3 {
		t.Errorf("expected length=3, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestUpvarintEncode4(t *testing.T) {
	var val uint32 = 134217728
	expected := []byte{0xC0, 0x80, 0x80, 0x00}
	res, length, err := EncodeUpvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 4 {
		t.Errorf("expected length=4, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestUpvarintEncode5(t *testing.T) {
	var val uint32 = 4294967295
	expected := []byte{0x8F, 0xFF, 0xFF, 0xFF, 0x07F}
	res, length, err := EncodeUpvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 5 {
		t.Errorf("expected length=5, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestPvarintEncode1(t *testing.T) {
	var val int32 = -1
	expected := []byte{0x7F}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 1 {
		t.Errorf("expected length=1, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}

	val = -64
	expected = []byte{0x40}
	res, length, err = EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != 1 {
		t.Errorf("expected length=1, actual=%v", length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}
	if res[0] != expected[0] {
		t.Errorf("expected res=%#x, actual=%#x", expected, res)
	}
}

func TestPvarintEncode11(t *testing.T) {
	var val int32 = -65
	expected := []byte{0xFF, 0x3F}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}

func TestPvarintEncode2(t *testing.T) {
	var val int32 = -4097
	expected := []byte{0xDF, 0x7F}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}

func TestPvarintEncode3(t *testing.T) {
	var val int32 = -8193
	expected := []byte{0xFF, 0xBF, 0x7F}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}

func TestPvarintEncode4(t *testing.T) {
	var val int32 = -2097152
	expected := []byte{0xFF, 0x80, 0x80, 0x00}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}

func TestPvarintEncode5(t *testing.T) {
	var val int32 = -134217729
	expected := []byte{0xFF, 0xBF, 0xFF, 0xFF, 0x7F}
	// var vic = new(VarIntCodec)
	// res := make([]byte, 10)
	// err := vic.EncodeInt32(res, val)
	res, _, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}

func TestPvarintEncode51(t *testing.T) {
	var val int32 = -2147483648
	expected := []byte{0xF8, 0x80, 0x80, 0x80, 0x00}
	res, length, err := EncodePvarint(val)
	if err != nil {
		t.Errorf("expected err=nil, actual=%v", err)
	}
	if length != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), length)
	}
	if len(res) != len(expected) {
		t.Errorf("expected length=%v, actual=%v", len(expected), len(res))
	}

	for i, v := range res {
		if v != expected[i] {
			t.Errorf("expected res=%#x, actual=%#x", expected[i], v)
		}
	}
}
