package codec

import (
	"errors"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// PrimitiveTag represents the Tag of TLTV
type PrimitiveTag struct {
	raw byte
}

// SeqID 获取Key的顺序ID
func (t *PrimitiveTag) SeqID() byte {
	return t.raw
}

// NewPrimitiveTag create a PrimitivePacket Tag field
func NewPrimitiveTag(b byte) (p *PrimitiveTag, err error) {
	if b&utils.MSB == utils.MSB {
		return nil, errors.New("not a primitive node")
	}

	return &PrimitiveTag{raw: b}, nil
}

// PrimitiveType represents the value type of TLTV
type PrimitiveType byte

const (
	// TypeString type data
	TypeString PrimitiveType = 0x00
	// TypeVarint 是可变长度的整数类型
	TypeVarint = 0x01
	// TypeFloat is IEEE754 format as big-endian
	TypeFloat = 0x02
	// TypeBoolean is True OR false
	TypeBoolean = 0x03
	// TypeUUID is 128-bits fixed-length
	TypeUUID = 0x04
	// TypeBinary 二进制数据
	TypeBinary = 0x40
	// TypeNode represent a node, other 7 bits used represent tag id
	TypeNode = 0x80
)

// IsValid checks primitive type
func (y PrimitiveType) IsValid() error {
	switch y {
	case TypeString, TypeVarint, TypeFloat, TypeBoolean, TypeUUID, TypeBinary, TypeNode:
		return nil
	}
	return errors.New("Invalid PrimitiveType")
}
