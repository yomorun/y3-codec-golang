package y3

import (
	"errors"
	"fmt"
)

// NodeTag represents the Tag of TLTV
type NodeTag struct {
	raw byte
}

// SeqID 获取Key的顺序ID
func (t *NodeTag) SeqID() byte {
	return t.raw & DropMSB
}

func (t *NodeTag) String() string {
	return fmt.Sprintf("Tag: raw=%4b, SeqID=%v", t.raw, t.SeqID())
}

func newNodeTag(b byte) (p *NodeTag, err error) {
	// 最高位始终为1
	if b&MSB != MSB {
		return nil, errors.New("not a node packet")
	}

	return &NodeTag{raw: b}, nil
}

// NodePacket 以`TLV结构`进行数据描述, 是用户定义类型
type NodePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag NodeTag
	// length + raw buffer
	basePacket
	// NodePackets 存储 Node 类型
	NodePackets []NodePacket
	// PrimitivePackets 存储 Primitive 类型
	PrimitivePackets []PrimitivePacket
}
