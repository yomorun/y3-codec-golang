package y3

import (
	"testing"
)

// 测试一个简单的node
// 若要表示一个JSON的结构：{
//	'0x03': {
//     '0x01': -1,
//     '0x02':  1,
//  },
//	'0x04': {
//     '0x01': -1,
//  },
// }
// YoMo Codec should ->
// 0x83 (is a node, sequence id=3)
//   0x10 (node value length is 8 bytes)
//     0x01, 0x04, 0x01, 0x01 (varint: -1)
//     0x02, 0x04, 0x01, 0x02 (varint: 1)
// 0x84 (is a node, sequence id=4)
//   0x08 (node value length is 8 bytes)
//     0x01, 0x04, 0x01, 0x01 (varint: -1)
func TestSimpleNode(t *testing.T) {
	// buf := []byte{0x83, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x01, 0x02, 0x84, 0x08, 0x01, 0x04, 0x01, 0x01}
	buf := []byte{0x83, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x01, 0x02}
	res, endPos, err := ReadNode(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.PrimitivePackets) != 2 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.NodePackets), 1)
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

	// if n1.Tag.SeqID() != 0x03 || n1.Length() /*.base.Length()*/ != 8 || len(n1.basePacket.raw) != 8 {
	// 	t.Errorf("n1 actual = %v", n1)
	// 	t.Errorf("n1.Tag.SeqID() actual = %v", n1.Tag.SeqID())
	// }

	if endPos != 10 {
		t.Errorf("endPos actual = %v, and Expected = %v", endPos, 10)
	}
}
