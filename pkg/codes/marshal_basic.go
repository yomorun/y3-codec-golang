package codes

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	y3 "github.com/yomorun/yomo-codec-golang"
)

var (
	startingToken byte = 0x01
)

func encodeBasic(observe byte, input interface{}) (buf []byte, err error) {
	encoder, err := marshalBasic(observe, input, nil)
	if err != nil {
		return []byte{}, err
	}
	return encoder.Encode(), nil
}

func marshalBasic(observe byte, input interface{}, root *y3.NodePacketEncoder) (encoder *y3.NodePacketEncoder, err error) {
	if observe == 0 {
		panic(fmt.Errorf("observe cannot be 0"))
	}

	if root == nil {
		root = y3.NewNodePacketEncoder(int(startingToken))
	}

	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.String:
		marshalBasicString(observe, input, root)
	case reflect.Int32:
		marshalBasicInt32(observe, input, root)
	case reflect.Uint32:
		marshalBasicUint32(observe, input, root)
	case reflect.Int64:
		marshalBasicInt64(observe, input, root)
	case reflect.Uint64:
		marshalBasicUint64(observe, input, root)
	case reflect.Float32:
		marshalBasicFloat32(observe, input, root)
	case reflect.Float64:
		marshalBasicFloat64(observe, input, root)
	case reflect.Array, reflect.Slice:
		marshalBasicSlice(observe, value, root)
	default:
		panic(errors.New("::Marshal error: no matching type"))
	}

	return root, nil
}

func marshalBasicSlice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	if value.Len() == 0 {
		return
	}

	switch value.Index(0).Kind() {
	case reflect.String:
		marshalBasicStringSlice(observe, value, encoder)
	case reflect.Int32:
		marshalBasicInt32Slice(observe, value, encoder)
	case reflect.Uint32:
		marshalBasicUint32Slice(observe, value, encoder)
	case reflect.Int64:
		marshalBasicInt64Slice(observe, value, encoder)
	case reflect.Uint64:
		marshalBasicUint64Slice(observe, value, encoder)
	case reflect.Float32:
		marshalBasicFloat32Slice(observe, value, encoder)
	case reflect.Float64:
		marshalBasicFloat64Slice(observe, value, encoder)
	default:
		panic(errors.New("::Marshal error: no matching type in Slice"))
	}
}

func marshalBasicString(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetStringValue(fmt.Sprintf("%v", input))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicInt32(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetInt32Value(input.(int32))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicUint32(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetUInt32Value(input.(uint32))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicInt64(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetInt64Value(input.(int64))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicUint64(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetUInt64Value(input.(uint64))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicFloat32(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetFloat32Value(input.(float32))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicFloat64(observe byte, input interface{}, encoder *y3.NodePacketEncoder) {
	var item = y3.NewPrimitivePacketEncoder(int(observe))
	item.SetFloat64Value(input.(float64))
	encoder.AddPrimitivePacket(item)
}

func marshalBasicStringSlice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToStringSliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicInt32Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicUint32Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt32Value(uint32(v.(uint64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicInt64Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicUint64Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt64Value(v.(uint64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicFloat32Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat32Value(float32(v.(float64)))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}

func marshalBasicFloat64Slice(observe byte, value reflect.Value, encoder *y3.NodePacketEncoder) {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat64Value(v.(float64))
			node.AddPrimitivePacket(item)
		}
	}
	encoder.AddNodePacket(node)
}
