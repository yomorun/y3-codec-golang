package spec

import (
	"testing"
)

func TestV2Nil(t *testing.T) {
	p, _ := NewPacket(18446744073709551487)
	p.SetNil()

	buf := compareBytes(t, p, []byte{0xFE, 0x7F, 0x00})

	dp, _ := FromBytes(buf)
	if dp.GetTag() != 18446744073709551487 {
		t.Errorf("result.Tag=% X, expect.Tag=% X", dp.GetTag(), uint64(18446744073709551487))
	}

	if dp.Length != 0 {
		t.Errorf("result.Length=%d, expect.Length=%v", dp.Length, 0)
	}

	if len(dp.valbuf) != 0 {
		t.Errorf("result.len(valbuf)=%d, expect.len(valbuf)=%v", len(dp.valbuf), 0)
	}
}

func TestV2UInt32(t *testing.T) {
	testV2UInt32(t, 0x03, 6, []byte{0x03, 0x01, 0x06})
	testV2UInt32(t, 0x06, 127, []byte{0x06, 0x02, 0x80, 0x7F})
}

func TestV2Int32(t *testing.T) {
	testV2Int32(t, 0x03, -1, []byte{0x03, 0x01, 0x7F})
	testV2Int32(t, 0x06, -65, []byte{0x06, 0x02, 0xFF, 0x3F})
	testV2Int32(t, 0x09, 255, []byte{0x09, 0x02, 0x81, 0x7F})
}

func TestV2UInt64(t *testing.T) {
	testV2UInt64(t, 0x03, 0, []byte{0x03, 0x01, 0x00})
	testV2UInt64(t, 0x06, 1, []byte{0x06, 0x01, 0x01})
	testV2UInt64(t, 0x09, 18446744073709551615, []byte{0x09, 0x01, 0x7F})
}

func TestV2Int64(t *testing.T) {
	testV2Int64(t, 0x03, 0, []byte{0x03, 0x01, 0x00})
	testV2Int64(t, 0x06, 1, []byte{0x06, 0x01, 0x01})
	testV2Int64(t, 0x09, -1, []byte{0x09, 0x01, 0x7F})
}

func TestV2Float32(t *testing.T) {
	testV2Float32(t, 0x03, -2, []byte{0x03, 0x01, 0xC0})
	testV2Float32(t, 0x06, 0.25, []byte{0x06, 0x02, 0x3E, 0x80})
	testV2Float32(t, 0x09, 68.123, []byte{0x09, 0x04, 0x42, 0x88, 0x3E, 0xFA})
}

func TestV2Float64(t *testing.T) {
	testV2Float64(t, 0x03, 23, []byte{0x03, 0x02, 0x40, 0x37})
	testV2Float64(t, 0x06, 2, []byte{0x06, 0x01, 0x40})
	testV2Float64(t, 0x09, 0.01171875, []byte{0x09, 0x02, 0x3F, 0x88})
}

func TestV2String(t *testing.T) {
	p, _ := NewPacket(0x01)
	p.SetUTF8String("C")
	compareBytes(t, p, []byte{0x01, 0x01, 0x43})
	p.SetUTF8String("CC")
	compareBytes(t, p, []byte{0x01, 0x02, 0x43, 0x43})
	p.SetUTF8String("Yona")
	compareBytes(t, p, []byte{0x01, 0x04, 0x59, 0x6F, 0x6E, 0x61})
	p.SetUTF8String("https://yomo.run")
	buf := compareBytes(t, p, []byte{0x01, 0x10, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3A, 0x2F, 0x2F, 0x79, 0x6F, 0x6D, 0x6F, 0x2E, 0x72, 0x75, 0x6E})

	dp, _ := FromBytes(buf)
	resval := dp.GetValueAsUTF8String()
	if resval != "https://yomo.run" {
		t.Errorf("result=%s, expect=https://yomo.run", resval)
	}
}

func TestV2Bytes(t *testing.T) {
	p, _ := NewPacket(0x01)
	p.PutBytes([]byte{0x03, 0x06, 0x09, 0x0C, 0x0F})
	compareBytes(t, p, []byte{0x01, 0x05, 0x03, 0x06, 0x09, 0x0C, 0x0F})
	p.PutBytes([]byte{0x06, 0x01, 0x01})
	compareBytes(t, p, []byte{0x01, 0x08, 0x03, 0x06, 0x09, 0x0C, 0x0F, 0x06, 0x01, 0x01})
}

func TestV2Bool(t *testing.T) {
	p, _ := NewPacket(0x01)
	p.SetBool(true)
	compareBytes(t, p, []byte{0x01, 0x01, 0x01})

	p.SetBool(false)
	buf := compareBytes(t, p, []byte{0x01, 0x01, 0x00})

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsBool()
	if resval != false {
		t.Errorf("result=%v, expect=%v", resval, false)
	}
}

func testV2UInt32(t *testing.T, tag uint64, val uint32, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetUInt32(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsUInt32()
	if resval != val {
		t.Errorf("result=%d, expect=%d", resval, val)
	}
}

func testV2Int32(t *testing.T, tag uint64, val int, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetInt32(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsInt32()
	if resval != int32(val) {
		t.Errorf("result=%d, expect=%d", resval, val)
	}
}

func testV2UInt64(t *testing.T, tag uint64, val uint64, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetUInt64(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsUInt64()
	if resval != val {
		t.Errorf("result=%d, expect=%d", resval, val)
	}
}

func testV2Int64(t *testing.T, tag uint64, val int64, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetInt64(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsInt64()
	if resval != val {
		t.Errorf("result=%d, expect=%d", resval, val)
	}
}

func testV2Float32(t *testing.T, tag uint64, val float32, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetFloat32(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsFloat32()
	if resval != val {
		t.Errorf("result=%f, expect=%f", resval, val)
	}
}

func testV2Float64(t *testing.T, tag uint64, val float64, expected []byte) {
	p, _ := NewPacket(tag)
	p.SetFloat64(val)

	buf := compareBytes(t, p, expected)

	dp, _ := FromBytes(buf)
	resval, _ := dp.GetValueAsFloat64()
	if resval != val {
		t.Errorf("result=%f, expect=%f", resval, val)
	}
}
