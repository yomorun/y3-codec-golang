package codes

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	y3 "github.com/yomorun/yomo-codec-golang"
)

// BasicDecoder: for UnmarshalBasic
type BasicDecoder struct {
	Observe byte
}

func newBasicDecoder(observe byte) *BasicDecoder {
	return &BasicDecoder{Observe: observe}
}

func (d BasicDecoder) Unmarshal(data []byte, mold *interface{}) error {
	nodePacket, _, err := y3.DecodeNodePacket(data)
	if err != nil {
		return err
	}

	return d.UnmarshalByNodePacket(nodePacket, mold)
}

func (d BasicDecoder) UnmarshalNative(data []byte, mold *interface{}) error {
	switch reflect.TypeOf(*mold).Kind() {
	case reflect.Array, reflect.Slice:
		nodePacket, _, err := y3.DecodeNodePacket(data)
		if err != nil {
			return err
		}
		if nodePacket.IsArray() && len(nodePacket.PrimitivePackets) > 0 {
			return d.unmarshalPrimitivePacketArray(nodePacket.PrimitivePackets, mold)
		}
		return fmt.Errorf("not be a packet that can be resolved. mold=%v", mold)
	default:
		primitivePacket, _, _, err := y3.DecodePrimitivePacket(data)
		if err != nil {
			return err
		}
		return d.unmarshalPrimitivePacket(*primitivePacket, mold)
	}
}

func (d BasicDecoder) UnmarshalByNodePacket(node *y3.NodePacket, mold *interface{}) error {
	key := d.Observe
	ok, isNode, packet := packetutils.MatchingKey(key, node)
	if !ok {
		return errors.New(fmt.Sprintf("not found mold in result. key:%#x", key))
	}

	return d.unmarshalPacket(packet, isNode, mold)
}

func (d BasicDecoder) unmarshalPacket(packet interface{}, isNode bool, mold *interface{}) error {
	if isNode == false {
		primitivePacket := packet.(y3.PrimitivePacket)
		return d.unmarshalPrimitivePacket(primitivePacket, mold)
	}

	nodePacket := packet.(y3.NodePacket)
	if nodePacket.IsArray() && len(nodePacket.PrimitivePackets) > 0 {
		return d.unmarshalPrimitivePacketArray(nodePacket.PrimitivePackets, mold)
	}

	return nil
}

// unmarshalPrimitivePacket convert PrimitivePacket to Mold
func (d BasicDecoder) unmarshalPrimitivePacket(primitivePacket y3.PrimitivePacket, mold *interface{}) error {
	switch reflect.TypeOf(*mold).Kind() {
	case reflect.String:
		v, err := primitivePacket.ToUTF8String()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Int32:
		v, err := primitivePacket.ToInt32()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Uint32:
		v, err := primitivePacket.ToUInt32()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Int64:
		v, err := primitivePacket.ToInt64()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Uint64:
		v, err := primitivePacket.ToUInt64()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Float32:
		v, err := primitivePacket.ToFloat32()
		if err != nil {
			return err
		}
		*mold = v
	case reflect.Float64:
		v, err := primitivePacket.ToFloat64()
		if err != nil {
			return err
		}
		*mold = v
	}

	return nil
}

// unmarshalPrimitivePacketArray convert []PrimitivePacket to Mold
func (d BasicDecoder) unmarshalPrimitivePacketArray(primitivePackets []y3.PrimitivePacket, mold *interface{}) error {
	result := make([]interface{}, 0)
	switch reflect.TypeOf(*mold).Kind() {
	case reflect.Array, reflect.Slice:
		for _, p := range primitivePackets {
			switch reflect.TypeOf(*mold).Elem().Kind() {
			case reflect.String:
				v, _ := p.ToUTF8String()
				result = append(result, v)
			case reflect.Int32:
				v, _ := p.ToInt32()
				result = append(result, v)
			case reflect.Uint32:
				v, _ := p.ToUInt32()
				result = append(result, v)
			case reflect.Int64:
				v, _ := p.ToInt64()
				result = append(result, v)
			case reflect.Uint64:
				v, _ := p.ToUInt64()
				result = append(result, v)
			case reflect.Float32:
				v, _ := p.ToFloat32()
				result = append(result, v)
			case reflect.Float64:
				v, _ := p.ToFloat64()
				result = append(result, v)
			}
		}
	}
	*mold = result
	return nil
}