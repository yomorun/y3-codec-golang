package y3

import (
	"fmt"

	codec "github.com/yomorun/yomo-codec-golang/internal/codec"
)

// 描述最小的Packet大小为4个字节
const primitivePacketBufferMinimalLength = 4

// PrimitivePacket 定义了值类型的节点，是Codec中的最小单位，以`TLTV结构`进行数据描述
type PrimitivePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag byte
	// length and raw buffer in base
	basePacket
	// 描述Value的数据类型
	Type codec.PrimitiveType
}

// ToInt64 parse raw as int64 value
func (p *PrimitivePacket) ToInt64() (int64, error) {
	dec, _ := codec.NewDecoder(p.basePacket.raw)
	result, err := dec.Decode()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// ToUTF8String parse raw data as string value
func (p *PrimitivePacket) ToUTF8String() (string, error) {
	return string(p.basePacket.raw), nil
}

// 用于打印时使用
func (p *PrimitivePacket) String() string {
	return fmt.Sprintf("Tag=%v, Length=%v, Type=%v, RawDataLength=%v, Raw=[% x]", p.Tag, p.Length(), p.Type, len(p.basePacket.raw), p.basePacket.raw)
}
