package y3

import (
	"testing"
)

// 测试一个简单的node
// 若要表示一个JSON的结构：{
// '0x04': {
//   '0x01': -1,
// },
// YoMo Codec should ->
// 0x84 (is a node, sequence id=4)
//   0x08 (node value length is 4 bytes)
//     0x01, 0x04, 0x01, 0x01 (varint: -1)
func TestSimple1Node(t *testing.T) {
	buf := []byte{0x84, 0x08, 0x01, 0x04, 0x01, 0x01}
	res, endPos, err := ReadNode(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.PrimitivePackets) != 1 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.NodePackets), 1)
	}

	v1, err := res.PrimitivePackets[0].ToInt64()
	if err != nil {
		t.Error(err)
	}

	if v1 != int64(-1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v1, -1)
	}

	if endPos != 6 {
		t.Errorf("endPos actual = %v, and Expected = %v", endPos, 10)
	}
}

// 测试一个简单的node
// 若要表示一个JSON的结构：
// '0x03': {
//   '0x01': -1,
//   '0x02':  1,
// },
// YoMo Codec should ->
// 0x83 (is a node, sequence id=3)
//   0x10 (node value length is 8 bytes)
//     0x01, 0x04, 0x01, 0x01 (varint: -1)
//     0x02, 0x04, 0x01, 0x02 (varint: 1)
func TestSimple2Nodes(t *testing.T) {
	buf := []byte{0x83, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x01, 0x02}
	res, endPos, err := ReadNode(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.PrimitivePackets) != 2 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.NodePackets), 2)
	}

	v1, err := res.PrimitivePackets[0].ToInt64()
	if err != nil {
		t.Error(err)
	}

	if v1 != int64(-1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v1, -1)
	}

	v2, err := res.PrimitivePackets[1].ToInt64()
	if err != nil {
		t.Error(err)
	}

	if v2 != int64(1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v2, 1)
	}

	if endPos != 10 {
		t.Errorf("endPos actual = %v, and Expected = %v", endPos, 10)
	}
}

// 测试一个简单的node
// 若要表示一个JSON的结构：
// '0x05': {
//	'0x04': {
//     '0x01': -1,
//     '0x02':  1,
//  },
//	'0x04': {
//     '0x01': -1,
//  },
// }
// YoMo Codec should ->
// 0x85
//   0x20(node value length is 16 bytes)
//     0x84 (is a node, sequence id=3)
//       0x10 (node value length is 8 bytes)
//         0x01, 0x04, 0x01, 0x01 (varint: -1)
//         0x02, 0x04, 0x00, 0x43 (varint: 1)
//     0x83 (is a node, sequence id=4)
//       0x08 (node value length is 4 bytes)
//         0x01, 0x04, 0x01, 0x03 (varint: -2)
func TestComplexNodes(t *testing.T) {
	buf := []byte{0x85, 0x20, 0x84, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x00, 0x43, 0x83, 0x08, 0x01, 0x04, 0x01, 0x03}
	res, endPos, err := ReadNode(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if endPos != len(buf) {
		t.Errorf("endPos actual = %v, and expected = %v", endPos, len(buf))
	}

	if len(res.NodePackets) != 2 {
		t.Errorf("res.NodePackets actual = %v, and expected = %v", len(res.NodePackets), 2)
	}

	if len(res.PrimitivePackets) != 0 {
		t.Errorf("res.PrimitivePackets actual = %v, and expected = %v", len(res.PrimitivePackets), 0)
	}

	n1 := res.NodePackets[0]
	if len(n1.PrimitivePackets) != 2 {
		t.Errorf("n1.PrimitivePackets actual = %v, and expected = %v", len(n1.PrimitivePackets), 2)
	}

	n1p1, _ := n1.PrimitivePackets[0].ToInt64()
	n1p2, _ := n1.PrimitivePackets[1].ToUTF8String()

	n2 := res.NodePackets[1]
	if len(n2.PrimitivePackets) != 1 {
		t.Errorf("n2.PrimitivePackets actual = %v, and expected = %v", len(n2.PrimitivePackets), 1)
	}

	n2p1, _ := n2.PrimitivePackets[0].ToInt64()

	if n1p1 != -1 || n1p2 != "C" || n2p1 != -2 {
		t.Errorf("n1p1=%v, n1p2=%v, n2p1=%v", n1p1, n1p2, n2p1)
	}

	// if n1.Tag.SeqID() != 0x03 || n1.Length() /*.base.Length()*/ != 8 || len(n1.basePacket.raw) != 8 {
	// 	t.Errorf("n1 actual = %v", n1)
	// 	t.Errorf("n1.Tag.SeqID() actual = %v", n1.Tag.SeqID())
	// }
}
