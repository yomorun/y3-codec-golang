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
	buf := []byte{0x83, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x01, 0x02, 0x84, 0x08, 0x01, 0x04, 0x01, 0x01}
	res, err := ReadAll(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.nodes) != 2 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.nodes), 2)
	}

	n1 := res.nodes[0]
	if n1.Tag.SeqID() != 0x03 || n1.Length != 8 || len(n1.Raw) != 8 {
		t.Errorf("n1 actual = %v", n1)
		t.Errorf("n1.Tag.SeqID() actual = %v", n1.Tag.SeqID())
		// t.Errorf("data parse error, Length=%v|%v", n1.Length, n2.Length)
		// t.Errorf("data parse error, len(Raw)=%v|%v", len(n1.Raw), len(n2.Raw))
	}

	n2 := res.nodes[1]
	if n2.Tag.SeqID() != 0x04 || n2.Length != 4 || len(n2.Raw) != 4 {
		t.Errorf("n2 actual = %v", n2)
	}
}
