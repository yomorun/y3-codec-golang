package y3

import (
	"fmt"

	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"
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
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt32(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUInt32 parse raw as int32 value
func (p *PrimitivePacket) ToUInt32() (uint32, error) {
	var val uint32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt32(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToInt64 parse raw as int32 value
func (p *PrimitivePacket) ToInt64() (int64, error) {
	var val int64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt64(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUInt64 parse raw as uint64 value
func (p *PrimitivePacket) ToUInt64() (uint64, error) {
	var val uint64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt64(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToFloat32 parse raw as float32 value
func (p *PrimitivePacket) ToFloat32() (float32, error) {
	var val float32
	codec := encoding.VarCodec{Size: len(p.valbuf)}
	err := codec.DecodeVarFloat32(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToFloat64 parse raw as float64 value
func (p *PrimitivePacket) ToFloat64() (float64, error) {
	var val float64
	codec := encoding.VarCodec{Size: len(p.valbuf)}
	err := codec.DecodeVarFloat64(p.valbuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToBool parse raw as bool value
func (p *PrimitivePacket) ToBool() (bool, error) {
	var val bool
	codec := encoding.VarCodec{Size: len(p.valbuf)}
	err := codec.DecodePVarBool(p.valbuf, &val)
	if err != nil {
		return false, err
	}
	return val, nil
}

// ToUTF8String parse raw data as string value
func (p *PrimitivePacket) ToUTF8String() (string, error) {
	return string(p.valbuf), nil
}

//// HasPacketArray determine if the value is an array
//func (p *PrimitivePacket) HasPacketArray() bool {
//	return p.tag.IsArray()
//}

//// ToPacketArray parse raw data as primitive packet array
//func (p *PrimitivePacket) ToPacketArray() (arr []*PrimitivePacket, err error) {
//	arr = make([]*PrimitivePacket, 0)
//	if !p.HasPacketArray() {
//		return arr, errors.New("value is not an array")
//	}
//
//	buf := p.valbuf
//
//	for {
//		packet, _, size, err := DecodePrimitivePacket(buf)
//		if err != nil {
//			return arr, err
//		}
//
//		arr = append(arr, packet)
//
//		tlvLength := 1 + uint32(size) + packet.length
//		if uint32(len(buf)) > tlvLength {
//			buf = buf[tlvLength:]
//			continue
//		}
//		break
//	}
//
//	return arr, nil
//}

func (p *PrimitivePacket) ToBytes() []byte {
	return p.valbuf
}
