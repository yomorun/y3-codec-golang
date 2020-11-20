package codes

import (
	"errors"
	"fmt"
	"reflect"

	y3 "github.com/yomorun/yomo-codec-golang"
)

func (codec *yomoCodec) UnmarshalBasic(data []byte, mold *interface{}) error {
	decoder := newBasicDecoder(codec.Observe)
	return decoder.Unmarshal(data, mold)
}

type BasicDecoder struct {
	Observe string
}

func newBasicDecoder(observe string) *BasicDecoder {
	return &BasicDecoder{Observe: observe}
}

func (d BasicDecoder) Unmarshal(data []byte, mold *interface{}) error {
	key := keyOf(d.Observe)
	pct, _, err := y3.DecodeNodePacket(data)
	if err != nil {
		return err
	}

	ok, isNode, packet := matchingKey(key, pct)
	if !ok {
		return errors.New(fmt.Sprintf("not found mold in result. key:%#x", key))
	}

	return d.unmarshalPrimitive(packet, isNode, mold)
}

func (d BasicDecoder) unmarshalPrimitive(packet interface{}, isNode bool, mold *interface{}) error {
	if isNode == false {
		primitivePacket := packet.(y3.PrimitivePacket)
		return d.unmarshalPrimitivePacket(primitivePacket, mold)
	}

	//fmt.Printf("#78 reflect.TypeOf(*output).Kind()=%v\n", reflect.TypeOf(*output).Kind())
	nodePacket := packet.(y3.NodePacket)
	if nodePacket.IsArray() && len(nodePacket.PrimitivePackets) > 0 {
		return d.unmarshalPrimitivePacketArray(nodePacket.PrimitivePackets, mold)
	}

	return nil
}

// convertPrimitivePacketToMold convert PrimitivePacket to Mold
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

	// TODO: unfinished: conv to custom struct

	return nil
}

// convertPrimitivePacketArrayToMold convert []PrimitivePacket to Mold
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
