package y3

import (
	"errors"
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

// Val 是TLTV中的Value
type Val struct {
	raw []byte
}

// BasePacket 是YoMo Codec中最小的单元，以`TLTV结构`进行数据描述
type BasePacket struct {
	// Tag 是TLTV中的Tag, 描述Key
	Tag byte
	// ValueType + Value 的字节长度
	Length int64
	// 描述Value的数据类型
	Type Type
	// Value的字节
	Val *Val
	// Raw bytes
	raw []byte
}

// PacketBufferMinimalLength 描述最小的Packet大小为4个字节
const PacketBufferMinimalLength = 4
