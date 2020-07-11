package y3

import (
	"errors"
	"fmt"

	varint "github.com/yomorun/yomo-codec-golang/internal"
)

// Type represents the value type of TLTV
type Type byte

const (
	// String type data
	String Type = 0x00
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
	// Node represent a node
	Node = 0x80
)

func (y Type) isValid() error {
	switch y {
	case String, Varint, Float, Boolean, UUID, Binary, Node:
		return nil
	}
	return errors.New("Invalid Type")
}

// BasePacket 是YoMo Codec中最小的单元，以`TLTV结构`进行数据描述, 解析出primitive type的value
type BasePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag byte
	// ValueType + Value 的字节长度
	Length int64
	// 描述Value的数据类型
	Type Type
	// Value的字节
	Val interface{}
	// Raw bytes
	raw []byte
}

// ToInt64 parse raw as int64 value
func (p *BasePacket) ToInt64() (int64, error) {
	dec, _ := varint.NewDecoder(p.raw)
	result, err := dec.Decode()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (p *BasePacket) String() string {
	return fmt.Sprintf("Tag=%v, Length=%v, Type=%v, RawDataLength=%v", p.Tag, p.Length, p.Type, len(p.raw))
}

// PacketBufferMinimalLength 描述最小的Packet大小为4个字节
const PacketBufferMinimalLength = 4
