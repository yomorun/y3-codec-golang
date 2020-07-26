package y3

import (
	"fmt"

	encoding "github.com/yomorun/yomo-codec-golang/pkg"
)

// 描述最小的Packet大小为4个字节
const primitivePacketBufferMinimalLength = 3

// PrimitivePacket 定义了值类型的节点，是Codec中的最小单位，以`TLV结构`进行数据描述
type PrimitivePacket basePacket

// SeqID returns the key of primitive packet
func (p *PrimitivePacket) SeqID() byte {
	return p.tag.SeqID()
}

// String prints debug info
func (p *PrimitivePacket) String() string {
	return fmt.Sprintf("Tag=%#x, Length=%v, RawDataLength=%v, Raw=[%#x]", p.tag, p.length, len(p.valbuf), p.valbuf)
}

// ToInt32 parse raw as int32 value
func (p *PrimitivePacket) ToInt32() (int32, error) {
	var val int32
	codec := encoding.VarIntCodec{}
	err := codec.DecodePVarInt32(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUInt32 parse raw as int32 value
func (p *PrimitivePacket) ToUInt32() (uint32, error) {
	var val uint32
	codec := encoding.VarIntCodec{}
	err := codec.DecodePVarUInt32(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUTF8String parse raw data as string value
func (p *PrimitivePacket) ToUTF8String() (string, error) {
	return string(p.valbuf), nil
}
