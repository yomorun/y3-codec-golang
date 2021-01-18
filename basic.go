package y3

import (
	"fmt"
	"reflect"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

type BasicEncoder interface {
	Encode(input interface{}) (buf []byte, err error)
}

func NewBasicEncoder(observe byte) BasicEncoder {
	return &basicEncoder{observe: observe}
}

type basicEncoder struct {
	observe byte
}

func (e basicEncoder) Encode(input interface{}) (buf []byte, err error) {
	encoder, err := e.marshalBasic(input, nil)
	if err != nil {
		return []byte{}, err
	}
	return encoder.Encode(), nil
}

func (e basicEncoder) marshalBasic(input interface{}, root *NodePacketEncoder) (encoder *NodePacketEncoder, err error) {
	if e.observe == 0 {
		panic(fmt.Errorf("observe cannot be 0"))
	}

	if root == nil {
		root = NewNodePacketEncoder(int(startingToken))
	}

	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.String:
		e.marshalBasicString(input, root)
	case reflect.Int32:
		e.marshalBasicInt32(input, root)
	case reflect.Uint32:
		e.marshalBasicUint32(input, root)
	case reflect.Int64:
		e.marshalBasicInt64(input, root)
	case reflect.Uint64:
		e.marshalBasicUint64(input, root)
	case reflect.Float32:
		e.marshalBasicFloat32(input, root)
	case reflect.Float64:
		e.marshalBasicFloat64(input, root)
	case reflect.Bool:
		e.marshalBasicBool(input, root)
	case reflect.Array, reflect.Slice:
		e.marshalBasicSlice(value, root)
	default:
		panic(fmt.Errorf("marshal error, no matching type: %v", value.Kind()))
	}

	return root, nil
}

func (e basicEncoder) marshalBasicSlice(value reflect.Value, encoder *NodePacketEncoder) {
	if value.Len() == 0 {
		return
	}

	switch value.Index(0).Kind() {
	case reflect.String:
		e.marshalBasicStringSlice(value, encoder)
	case reflect.Int32:
		e.marshalBasicInt32Slice(value, encoder)
	case reflect.Uint32:
		e.marshalBasicUint32Slice(value, encoder)
	case reflect.Int64:
		e.marshalBasicInt64Slice(value, encoder)
	case reflect.Uint64:
		e.marshalBasicUint64Slice(value, encoder)
	case reflect.Float32:
		e.marshalBasicFloat32Slice(value, encoder)
	case reflect.Float64:
		e.marshalBasicFloat64Slice(value, encoder)
	case reflect.Bool:
		e.marshalBasicBoolSlice(value, encoder)
	default:
		panic(fmt.Errorf("marshal error, no matching type in Slice: %v", value.Index(0).Kind()))
	}
}

func (e basicEncoder) marshalBasicString(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetStringValue(fmt.Sprintf("%v", input))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicInt32(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetInt32Value(input.(int32))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicUint32(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetUInt32Value(input.(uint32))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicInt64(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetInt64Value(input.(int64))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicUint64(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetUInt64Value(input.(uint64))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicFloat32(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetFloat32Value(input.(float32))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicFloat64(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetFloat64Value(input.(float64))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicBool(input interface{}, encoder *NodePacketEncoder) {
	var item = NewPrimitivePacketEncoder(int(e.observe))
	item.SetBoolValue(input.(bool))
	encoder.AddPrimitivePacket(item)
}

func (e basicEncoder) marshalBasicStringSlice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToStringSliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicInt32Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicUint32Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt32Value(uint32(v.(uint64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicInt64Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicUint64Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt64Value(v.(uint64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicFloat32Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat32Value(float32(v.(float64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicFloat64Slice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat64Value(v.(float64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func (e basicEncoder) marshalBasicBoolSlice(value reflect.Value, encoder *NodePacketEncoder) {
	var node = NewNodeArrayPacketEncoder(int(e.observe))
	if out, ok := utils.ToBoolSliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetBoolValue(v.(bool))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}
