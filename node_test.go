package y3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Assume a JSON object like this：
// '0x04': {
//   '0x01': -1,
// },
// YoMo Codec should ->
// 0x84 (is a node, sequence id=4)
//   0x03 (node value length is 4 bytes)
//     0x01, 0x01, 0x7F (pvarint: -1)
func TestSimple1Node(t *testing.T) {
	buf := []byte{0x84, 0x03, 0x01, 0x01, 0x7F}
	res, packetLength, err := DecodeNodePacket(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(buf), packetLength)
	assert.Equal(t, 0, len(res.NodePackets))
	assert.Equal(t, 1, len(res.PrimitivePackets))
	assert.EqualValues(t, 0x04, res.SeqID())

	v1, err := res.PrimitivePackets[0].ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(-1), v1)
	assert.Equal(t, 5, packetLength)
}

// Assume a JSON object like this：
// '0x03': {
//   '0x01': -1,
//   '0x02':  1,
// },
// YoMo Codec should ->
// 0x83 (is a node, sequence id=3)
//   0x06 (node value length is 8 bytes)
//     0x01, 0x01, 0x7F (pvarint: -1)
//     0x02, 0x01, 0x01 (pvarint: 1)
func TestSimple2Nodes(t *testing.T) {
	buf := []byte{0x83, 0x06, 0x01, 0x01, 0x7F, 0x02, 0x01, 0x01}
	res, packetLength, err := DecodeNodePacket(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(buf), packetLength)
	assert.Equal(t, 0, len(res.NodePackets))
	assert.Equal(t, 2, len(res.PrimitivePackets))

	v1, err := res.PrimitivePackets[0].ToInt32()
	assert.NoError(t, err)
	assert.EqualValues(t, -1, v1)

	v2, err := res.PrimitivePackets[1].ToInt32()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v2)
}

// Assume a JSON object like this：
// '0x05': {
//	'0x04': {
//     '0x01': -1,
//     '0x02':  1,
//  },
//	'0x03': {
//     '0x01': -2,
//  },
// }
// YoMo Codec should ->
// 0x85
//   0x0D(node value length is 16 bytes)
//     0x84 (is a node, sequence id=3)
//       0x06 (node value length is 8 bytes)
//         0x01, 0x01, 0x7F (varint: -1)
//         0x02, 0x01, 0x43 (string: "C")
//     0x83 (is a node, sequence id=4)
//       0x03 (node value length is 4 bytes)
//         0x01, 0x01, 0x7E (varint: -2)
func TestComplexNodes(t *testing.T) {
	buf := []byte{0x85, 0x0D, 0x84, 0x06, 0x01, 0x01, 0x7F, 0x02, 0x01, 0x43, 0x83, 0x03, 0x01, 0x01, 0x7E}
	res, packetLength, err := DecodeNodePacket(buf)
	assert.NoError(t, err)
	assert.Equal(t, packetLength, len(buf))
	assert.Equal(t, 2, len(res.NodePackets))
	assert.Equal(t, 0, len(res.PrimitivePackets))

	n1 := res.NodePackets[0]
	assert.Equal(t, 2, len(n1.PrimitivePackets))

	n1p1, _ := n1.PrimitivePackets[0].ToInt32()
	n1p2, _ := n1.PrimitivePackets[1].ToUTF8String()

	n2 := res.NodePackets[1]
	assert.Equal(t, 1, len(n2.PrimitivePackets))

	n2p1, _ := n2.PrimitivePackets[0].ToInt32()
	if n1p1 != -1 || n1p2 != "C" || n2p1 != -2 {
		t.Errorf("n1p1=%v, n1p2=%v, n2p1=%v", n1p1, n1p2, n2p1)
	}
}

func TestEmptyNode(t *testing.T) {
	buf := []byte{0x86, 0x00}
	res, packetLength, err := DecodeNodePacket(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(buf), packetLength)
	assert.Equal(t, 0, len(res.NodePackets))
	assert.Equal(t, 0, len(res.PrimitivePackets))
}

func TestSubEmptyNode(t *testing.T) {
	buf := []byte{0x86, 0x02, 0x83, 0x00}
	res, packetLength, err := DecodeNodePacket(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(buf), packetLength)
	assert.Equal(t, 1, len(res.NodePackets))
	assert.Equal(t, 0, len(res.PrimitivePackets))
	assert.Equal(t, 0, len(res.NodePackets[0].NodePackets))
	assert.Equal(t, 0, len(res.NodePackets[0].PrimitivePackets))
}
