package y3

import (
	"errors"
	"fmt"

	varint "github.com/yomorun/yomo-codec-golang/internal"
)

// MSB 描述了`1000 0000`, 用于表示后续字节仍然是该变长类型值的一部分
const MSB byte = 0x80

// DropMSB 描述了`0111 1111`, 用于去除标识位使用
const DropMSB = 0x7F

// PrimitivePacketBufferMinimalLength 描述最小的Packet大小为4个字节
const PrimitivePacketBufferMinimalLength = 4

// PrimitiveTag represents the Tag of TLTV
type PrimitiveTag struct {
	raw byte
}

// SeqID 获取Key的顺序ID
func (t *PrimitiveTag) SeqID() byte {
	return t.raw
}

func newPrimitiveTag(b byte) (p *PrimitiveTag, err error) {
	if b&MSB == MSB {
		return nil, errors.New("not a primitive node")
	}

	return &PrimitiveTag{raw: b}, nil
}

// PrimitiveType represents the value type of TLTV
type PrimitiveType byte

const (
	// String type data
	String PrimitiveType = 0x00
	// Varint 是可变长度的整数类型
	Varint = 0x01
	// Float is IEEE754 format as big-endian
	Float = 0x02
	// Boolean is True OR false
	Boolean = 0x03
	// UUID is 128-bits fixed-length
	UUID = 0x04
	// Binary 二进制数据
	Binary = 0x40
	// Node represent a node, other 7 bits used represent tag id
	Node = 0x80
)

func (y PrimitiveType) isValid() error {
	switch y {
	case String, Varint, Float, Boolean, UUID, Binary, Node:
		return nil
	}
	return errors.New("Invalid PrimitiveType")
}

// PrimitivePacket 定义了值类型的节点，是Codec中的最小单位，以`TLTV结构`进行数据描述
type PrimitivePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag byte
	// length and raw buffer in base
	basePacket
	// 描述Value的数据类型
	Type PrimitiveType
}

// ToInt64 parse raw as int64 value
func (p *PrimitivePacket) ToInt64() (int64, error) {
	dec, _ := varint.NewDecoder(p.basePacket.raw)
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
