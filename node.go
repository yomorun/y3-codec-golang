package y3

// NodePacket 以`TLV结构`进行数据描述, 是用户定义类型
type NodePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag byte
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
		n.Tag = b[pos]
		// fmt.Printf("pos=%v, n.Tag=%4b\n", pos, n.Tag)
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
