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

// ReadAll parse out whole buffer to
func ReadAll(b []byte) (*Nodes, error) {
	// fmt.Println(hex.Dump(b))
	if len(b) == 0 {
		return &Nodes{}, nil
	}

	nodeArr := make([]NodePacket, 0)

	pos := 0
	end := len(b)
	for {
		if pos >= end {
			break
		}
		// Tag is 1 byte
		n := new(NodePacket)
		tag, err := newNodeTag(b[pos])
		if err != nil {
			return nil, err
		}
		n.Tag = *tag
		// fmt.Printf("pos=%d, n.Tag=%v\n", pos, n.Tag)
		pos++
		// Length is `varint`
		len, bufLen, err := ParseVarintLength(b, pos)
		if err != nil {
			return nil, err
		}
		n.Length = len // Length的值是Type+Value的字节长度
		// fmt.Printf("pos=%v, n.Length=%v\n", pos, n.Length)
		pos += bufLen
		// Raw is n.Length length
		vl := int(len)
		n.Raw = make([]byte, vl)
		copy(n.Raw, b[pos:pos+vl])
		// fmt.Printf("pos=%v, n.Raw=%x\n", pos, n.Raw)
		pos += vl
		// commit
		nodeArr = append(nodeArr, *n)
	}

	res := &Nodes{nodeArr}
	return res, nil
}
