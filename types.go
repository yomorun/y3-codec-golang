package y3

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func ToObject(v []byte, output interface{}) error {
	output, err := NewStructDecoder(output).Decode(v)
	return err
}

func ToInt32(v []byte) (int32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToInt32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToInt32Slice(v []byte) ([]int32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]int32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt32()
		result = append(result, v)
	}
	return result, nil
}

func ToUInt32(v []byte) (uint32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToUInt32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToUInt32Slice(v []byte) ([]uint32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]uint32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUInt32()
		result = append(result, v)
	}
	return result, nil
}

func ToInt64(v []byte) (int64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToInt64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToInt64Slice(v []byte) ([]int64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]int64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt64()
		result = append(result, v)
	}
	return result, nil
}

func ToUInt64(v []byte) (uint64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToUInt64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToUInt64Slice(v []byte) ([]uint64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]uint64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUInt64()
		result = append(result, v)
	}
	return result, nil
}

func ToFloat32(v []byte) (float32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToFloat32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToFloat32Slice(v []byte) ([]float32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]float32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToFloat32()
		result = append(result, v)
	}
	return result, nil
}

func ToFloat64(v []byte) (float64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToFloat64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

func ToFloat64Slice(v []byte) ([]float64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]float64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToFloat64()
		result = append(result, v)
	}
	return result, nil
}

func ToBool(v []byte) (bool, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return false, err
	}
	value, err := primitivePacket.ToBool()
	if err != nil {
		return false, nil
	}
	return value, nil
}

func ToBoolSlice(v []byte) ([]bool, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]bool, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToBool()
		result = append(result, v)
	}
	return result, nil
}

func ToUTF8String(v []byte) (string, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return "", err
	}
	value, err := primitivePacket.ToUTF8String()
	if err != nil {
		return "", nil
	}
	return value, nil
}

func ToUTF8StringSlice(v []byte) ([]string, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsArray() || len(packet.PrimitivePackets) <= 0 {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]string, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUTF8String()
		result = append(result, v)
	}
	return result, nil
}
