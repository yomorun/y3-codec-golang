package y3

import (
	"bytes"
	"testing"

	"github.com/yomorun/y3-codec-golang/pkg/spec"
)

func TestStreamDecode1(t *testing.T) {
	// testStreamDecode(t, []byte{0x01, 0x00}, []byte{0x01}, []byte{0x00}, []byte{})
	data := []byte{0x01, 0x01, 0x03} //, []byte{0x01}, []byte{0x01}, []byte{0x03}, flag)
	// testStreamDecode(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05}, []byte{0x01}, []byte{0x02}, []byte{0x03, 0x04})
	// testStreamDecode(t, []byte{0x01, 0x03, 0x03, 0x04, 0x05}, []byte{0x01}, []byte{0x03}, []byte{0x03, 0x04, 0x05})
	// testStreamDecode(t, []byte{0x01, 0x01}, []byte{0x02}, []byte{0x01}, []byte{})
	// t.Errorf("---")

	expectTagbuf := []byte{0x01}
	expectLenbuf := []byte{0x01}
	expectValbuf := []byte{0x03}

	// as reader
	r := bytes.NewReader(data)
	// create steam decoder
	pr := NewStreamDecoder(r)

	// handler
	pr.OnPacket(func(p *spec.Packet) {
		t.Logf("[CALLBACK] p=%v", p)
		compareBytes(t, p.GetTagBuffer(), expectTagbuf, "T")
		compareBytes(t, p.GetLengthBuffer(), expectLenbuf, "L")
		compareBytes(t, p.GetValueBuffer(), expectValbuf, "V")
	})
	pr.Start()
}

func TestStreamDecode2(t *testing.T) {
	data := []byte{0x01, 0x00}

	expectTagbuf := []byte{0x01}
	expectLenbuf := []byte{0x00}
	expectValbuf := []byte{}

	// as reader
	r := bytes.NewReader(data)
	// create steam decoder
	pr := NewStreamDecoder(r)

	// handler
	pr.OnPacket(func(p *spec.Packet) {
		t.Logf("[CALLBACK] p=%v", p)
		compareBytes(t, p.GetTagBuffer(), expectTagbuf, "T")
		compareBytes(t, p.GetLengthBuffer(), expectLenbuf, "L")
		compareBytes(t, p.GetValueBuffer(), expectValbuf, "V")
	})
	pr.Start()
}

func TestStreamDecode3(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}

	expectTagbuf := []byte{0x01}
	expectLenbuf := []byte{0x02}
	expectValbuf := []byte{0x03, 0x04}

	// as reader
	r := bytes.NewReader(data)
	// create steam decoder
	pr := NewStreamDecoder(r)

	// handler
	pr.OnPacket(func(p *spec.Packet) {
		t.Logf("[CALLBACK] p=%v", p)
		compareBytes(t, p.GetTagBuffer(), expectTagbuf, "T")
		compareBytes(t, p.GetLengthBuffer(), expectLenbuf, "L")
		compareBytes(t, p.GetValueBuffer(), expectValbuf, "V")
	})
	pr.Start()
}

func TestStreamDecode4(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}

	// as reader
	r := bytes.NewReader(data)
	// create steam decoder
	pr := NewStreamDecoder(r)

	parsed := false

	// handler
	pr.OnPacket(func(p *spec.Packet) {
		t.Logf("[CALLBACK] p=%v", p)
		if p != nil {
			t.Error(p)
		}
		parsed = true
	})
	pr.Start()

	if parsed == true {
		t.Errorf("Should not trigger callback")
	}
}

func TestStreamDecode5(t *testing.T) {
	data := []byte{0x01, 0x01, 0x01, 0x02, 0x02, 0x01, 0x02, 0x03, 0x00, 0x04, 0x03, 0x01, 0x02, 0x03}

	// as reader
	r := bytes.NewReader(data)
	// create steam decoder
	pr := NewStreamDecoder(r)

	times := 1

	// handler
	pr.OnPacket(func(p *spec.Packet) {
		if times == 1 {
			compareBytes(t, p.GetTagBuffer(), []byte{0x01}, "T")
			compareBytes(t, p.GetLengthBuffer(), []byte{0x01}, "L")
			compareBytes(t, p.GetValueBuffer(), []byte{0x01}, "V")
		}
		if times == 2 {
			compareBytes(t, p.GetTagBuffer(), []byte{0x02}, "T")
			compareBytes(t, p.GetLengthBuffer(), []byte{0x02}, "L")
			compareBytes(t, p.GetValueBuffer(), []byte{0x01, 0x02}, "V")
		}
		if times == 3 {
			compareBytes(t, p.GetTagBuffer(), []byte{0x03}, "T")
			compareBytes(t, p.GetLengthBuffer(), []byte{0x00}, "L")
			compareBytes(t, p.GetValueBuffer(), []byte{}, "V")
			if p.Length != 0 {
				t.Errorf("Packet:Tag=[% X] Length should be 0, actual=%d", p.GetTagBuffer(), p.Length)
			}
		}
		if times == 4 {
			compareBytes(t, p.GetTagBuffer(), []byte{0x04}, "T")
			compareBytes(t, p.GetLengthBuffer(), []byte{0x03}, "L")
			compareBytes(t, p.GetValueBuffer(), []byte{0x01, 0x02, 0x03}, "V")
		}
		times++
	})

	pr.Start()
}

func compareBytes(t *testing.T, result []byte, expected []byte, v string) {
	if len(result) != len(expected) {
		t.Errorf("\n[%s] expected:[% X]\n actual:[% X]\n", v, expected, result)
	}

	for i, p := range result {
		if p != expected[i] {
			t.Errorf("\n[%s] expected:[% X]\n actual:[% X]\n", v, expected, result)
			break
		}
	}
}
