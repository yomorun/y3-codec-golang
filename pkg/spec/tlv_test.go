package spec

import (
	"testing"
)

func TestV2Tag(t *testing.T) {
	testV2Tags(t, 0x05, []byte{0x05})
	testV2Tags(t, 0x3F, []byte{0x3F})
	testV2Tags(t, 0x7F, []byte{0x80, 0x7F})
	testV2Tags(t, 0xFF, []byte{0x81, 0x7F})
	testV2Tags(t, 0xFFFF, []byte{0x83, 0xFF, 0x7F})
	testV2Tags(t, 0xFFFFFF, []byte{0x87, 0xFF, 0xFF, 0x7F})
}

func TestV2AddNode(t *testing.T) {
	parent, _ := NewPacket(0x01)

	child1, _ := NewPacket(0x02)
	child1.SetUInt32(3)

	parent.AddNode(child1)
	compareBytes(t, parent, []byte{0x01, 0x03, 0x02, 0x01, 0x03})

	child2, _ := NewPacket(0x03)
	child2.SetFloat64(0.01171875)

	parent.AddNode(child2)
	compareBytes(t, parent, []byte{0x01, 0x07, 0x02, 0x01, 0x03, 0x03, 0x02, 0x3F, 0x88})

	child3, _ := NewPacket(0x04)
	child3.SetFloat32(68.123)

	parent.AddNode(child3)
	compareBytes(t, parent, []byte{0x01, 0x0D, 0x02, 0x01, 0x03, 0x03, 0x02, 0x3F, 0x88, 0x04, 0x04, 0x42, 0x88, 0x3E, 0xFA})
}

func TestFromBytes(t *testing.T) {
	testFromBytes(t, []byte{0x03, 0x01, 0x02}, &Packet{idTag: 3, Length: 1, valbuf: []byte{0x02}})
	testFromBytes(t, []byte{0x81, 0x03, 0x01, 0x06}, &Packet{idTag: 131, Length: 1, valbuf: []byte{0x06}})

	foo := make([]byte, 129)
	bar := []byte{0x81, 0x03, 0x81, 0x01}
	buf := append(bar, foo...)
	testFromBytes(t, buf, &Packet{idTag: 131, Length: 129, valbuf: foo})
}

func TestReadTag(t *testing.T) {
	testReadTag(t, []byte{0xFF, 0x7F, 0x00}, 0, 18446744073709551615, 2)
	testReadTag(t, []byte{0x03, 0x01, 0x02}, 0, 3, 1)
	testReadTag(t, []byte{0xFF, 0x04, 0x01, 0x02}, 1, 4, 2)
	testReadTag(t, []byte{0x81, 0x05, 0x01, 0x02}, 0, 133, 2)
	testReadTag(t, []byte{0xFF, 0x81, 0x05, 0x01, 0x02}, 1, 133, 3)
}

func TestReadLength(t *testing.T) {
	testReadLength(t, []byte{0xFF, 0x7F, 0x00}, 2, 0, 3)
	testReadLength(t, []byte{0x03, 0x01, 0x02}, 1, 1, 2)
	testReadLength(t, []byte{0xFF, 0x04, 0x02, 0x00, 0x00}, 2, 2, 3)
	testReadLength(t, []byte{0x05, 0x81, 0x05, 0x00}, 1, 133, 3)
}

func testReadTag(t *testing.T, buffer []byte, position int, expectValue uint64, expectedCursor int) {
	var result uint64
	cursor, err := readPVarUInt64(buffer, position, &result)
	if err != nil {
		t.Error(err)
	}
	if result != expectValue {
		t.Errorf("expect tag=%d, actual=%d", expectValue, result)
	}
	if cursor != expectedCursor {
		t.Errorf("expect cursor=%d, actual=%d", expectedCursor, cursor)
	}
}

func testReadLength(t *testing.T, buffer []byte, position int, expectValue uint64, expectedCursor int) {
	var result uint64
	cursor, err := readPVarUInt64(buffer, position, &result)
	if err != nil {
		t.Errorf(">>> Got err=%s", err.Error())
	}
	if result != expectValue {
		t.Errorf("expect length=%d, actual=%d", expectValue, result)
	}
	if cursor != expectedCursor {
		t.Errorf("expect cursor=%d, actual=%d", expectedCursor, cursor)
	}
}

func testFromBytes(t *testing.T, buffer []byte, expected *Packet) {
	expected.Encode()
	res, err := FromBytes(buffer)
	// t.Logf("testFromBytes: res=%v, exptected=%v", res, expected)
	if err != nil {
		t.Error(err)
	} else {
		if res.GetTag() != expected.GetTag() {
			t.Errorf("Tag expected=[% X], actual=[% X]", expected.GetTag(), res.GetTag())
		}
		if res.Length != expected.Length {
			t.Errorf("Length expected=[% X], actual=[% X]", expected.Length, res.Length)
		}
		for i := range res.valbuf {
			if res.valbuf[i] != expected.valbuf[i] {
				t.Errorf("valbuf on [%d] expected=[% X], actual=[% X]", i, expected.valbuf[i], res.valbuf[i])
			}
		}
	}
}

func testV2Tags(t *testing.T, id uint64, expected []byte) {
	p, err := NewPacket(id)
	if err != nil {
		t.Errorf("TestV2Tag err=%v", err)
	}
	p.SetInt32(0)
	expected = append(expected, []byte{0x01, 0x00}...)
	compareBytes(t, p, expected)
}

func compareBytes(t *testing.T, p *Packet, expected []byte) []byte {
	result, err := p.Encode()
	if err != nil {
		t.Errorf("compareBytes error=%v", err)
	}
	for i, p := range result {
		if p != expected[i] {
			t.Errorf("\nexpected:[% X]\n  actual:[% X]\n", expected, result)
			break
		}
	}
	return result
}
