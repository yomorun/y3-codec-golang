package y3

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
