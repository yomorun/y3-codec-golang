package y3

import (
	"testing"
)

// JSON-like node:
// { '0x01': -1 }
// YoMo Codec should ->
// 0x01 (is a node, sequence id=1)
//   0x01 (node value length is 1 byte)
//     0x7F (pvarint: -1)
func TestEncoderPrimitiveInt64(t *testing.T) {
	expected := []byte{0x01, 0x01, 0x7F}
	var enc = CreateEncoder()

	// 0x01 - SeqID=1
	var prim = enc.CreatePrimitivePacket(0x01)
	// Value = -1
	prim.SetInt64Value(-1)

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}

// JSON-like node:
// { '0x01': "YoMo" }
// YoMo Codec should ->
// 0x01 (is a node, sequence id=1)
//   0x04 (pvarint, node value length is 4 bytes)
//     0x59, 0x6F, 0x4D, 0x6F (utf-8 string: "YoMo")
func TestEncoderPrimitiveString(t *testing.T) {
	expected := []byte{0x01, 0x04, 0x59, 0x6F, 0x4D, 0x6F}
	var enc = CreateEncoder()

	// 0x01 - SeqID=1
	var prim = enc.CreatePrimitivePacket(0x01)
	// Value = "YoMo"
	prim.SetStringValue("YoMo")

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}

// 0x81 : {
//   0x02: "YoMo",
// },
func TestEncoderNode1(t *testing.T) {
	expected := []byte{0x81, 0x06, 0x02, 0x04, 0x59, 0x6F, 0x4D, 0x6F}
	var enc = CreateEncoder()
	var node = enc.CreateNodePacket(0x01)
	var prim = enc.CreatePrimitivePacket(0x02)
	prim.SetStringValue("YoMo")
	node.AddPrimitivePacket(prim)
	enc.AddNodePacket(node)
	res := enc.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}
