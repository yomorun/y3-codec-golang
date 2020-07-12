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
	// ValueType + Value 的字节长度
	Length int64
	// Raw bytes
	Raw []byte
}

// Nodes 是解析后的整个Object
type Nodes struct {
	nodes []NodePacket
}
