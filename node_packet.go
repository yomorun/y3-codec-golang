package y3

// NodePacket 以`TLV结构`进行数据描述, 是用户定义类型
type NodePacket struct {
	// length + raw buffer
	basePacket
	// NodePackets 存储 Node 类型
	NodePackets []NodePacket
	// PrimitivePackets 存储 Primitive 类型
	PrimitivePackets []PrimitivePacket
}

// SeqID returns Tag Key
func (n *NodePacket) SeqID() byte {
	return n.basePacket.tag.SeqID()
}

// IsSlice determine if the current node is a Slice
func (n *NodePacket) IsSlice() bool {
	return n.basePacket.tag.IsSlice()
}

// GetValBuf get raw buffer of NodePacket
func (n *NodePacket) GetValBuf() []byte {
	return n.valBuf
}
