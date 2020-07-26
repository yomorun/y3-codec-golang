package y3

import (
	"testing"
)

// JSON-like node:
// { '0x01': -1 }
// YoMo Codec should ->
// 0x01 (is a node, sequence id=1)
//   0x01 (node value length is 1 byte)
//     0x01 (pvarint: -1)
func TestEncoderPrimitiveInt32(t *testing.T) {
	expected := []byte{0x01, 0x01, 0x7F}
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
	// Value = -1
	prim.SetInt32Value(-1)

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
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
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
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
	var prim = NewPrimitivePacketEncoder(0x02)
	prim.SetStringValue("YoMo")
	var node = NewNodePacketEncoder(0x01)
	node.AddPrimitivePacket(prim)
	res := node.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
		}
	}
}

// type bar struct {
// 	Name string
// }

// type foo struct {
// 	ID int
// 	*bar
// }
//
// var obj = &foo{ID: 1, bar: &bar{Name: "C"}}
//
// encode obj as:
//
// 0x81: {
//   0x02: 1,
//   0x83 : {
//     0x04: "C",
//   },
// }
//
// to
//
// [0x81, 0x08, 0x02, 0x01, 0x01, 0x83, 0x03, 0x04, 0x01, 0x43]
func TestEncoderNode2(t *testing.T) {
	expected := []byte{0x81, 0x08, 0x02, 0x01, 0x01, 0x83, 0x03, 0x04, 0x01, 0x43}
	// 0x81 - node
	var node1 = NewNodePacketEncoder(0x01)
	// 0x02 - ID=1
	var prim1 = NewPrimitivePacketEncoder(0x02)
	prim1.SetInt32Value(1)
	node1.AddPrimitivePacket(prim1)

	// 0x83 - &bar{}
	var node2 = NewNodePacketEncoder(0x03)

	// 0x04 - Name: "C"
	var prim2 = NewPrimitivePacketEncoder(0x04)
	prim2.SetStringValue("C")
	node2.AddPrimitivePacket(prim2)

	node1.AddNodePacket(node2)

	res := node1.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
		}
	}
}
