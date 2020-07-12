package y3

import codec "github.com/yomorun/yomo-codec-golang/internal/codec"

// NodePacket 以`TLV结构`进行数据描述, 是用户定义类型
type NodePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag codec.NodeTag
	// length + raw buffer
	basePacket
	// NodePackets 存储 Node 类型
	NodePackets []NodePacket
	// PrimitivePackets 存储 Primitive 类型
	PrimitivePackets []PrimitivePacket
}
