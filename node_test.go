package y3

import (
	"testing"
)

// 测试一个简单的node
// 若要表示一个JSON的结构：{
//	'0x03': {
//     '0x01': -1,
//     '0x02':  1,
//  }
//
// YoMo Codec should ->
// 0x83 (is a node)
//   0x10 (node value length is 8 bytes)
//     0x01, 0x04, 0x01, 0x01 (varint: -1)
//     0x02, 0x04, 0x01, 0x02 (varint: 1)
func TestSimpleNode(t *testing.T) {
	buf := []byte{0x81, 0x10, 0x01, 0x04, 0x01, 0x01, 0x02, 0x04, 0x01, 0x02}
	_, err := ReadAll(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}
}
