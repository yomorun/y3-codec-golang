package codes

import (
	"errors"
	"fmt"
	"reflect"

	y3 "github.com/yomorun/yomo-codec-golang"
	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// marshalPrimitive: marshal primitive to []byte
func marshalPrimitive(observe byte, input interface{}) (buf []byte, err error) {
	switch reflect.ValueOf(input).Kind() {
	case reflect.String:
		buf, err = marshalString(input)
	case reflect.Int32:
		buf, err = marshalInt32(input)
	case reflect.Uint32:
		buf, err = marshalUint32(input)
	case reflect.Int64:
		buf, err = marshalInt64(input)
	case reflect.Uint64:
		buf, err = marshalUint64(input)
	case reflect.Float32:
		buf, err = marshalFloat32(input)
	case reflect.Float64:
		buf, err = marshalFloat64(input)
	case reflect.Array, reflect.Slice:
		if reflect.ValueOf(input).Len() == 0 {
			break
		}
		switch reflect.ValueOf(input).Index(0).Elem().Kind() {
		case reflect.String:
			buf = marshalStringSlice(observe, input)
		case reflect.Int32:
			buf = marshalInt32Slice(observe, input)
		case reflect.Uint32:
			buf = marshalUint32Slice(observe, input)
		case reflect.Int64:
			buf = marshalInt64Slice(observe, input)
		case reflect.Uint64:
			buf = marshalUint64Slice(observe, input)
		case reflect.Float32:
			buf = marshalFloat32Slice(observe, input)
		case reflect.Float64:
			buf = marshalFloat64Slice(observe, input)
		default:
			panic(errors.New("::Marshal error: no matching type in Slice"))
		}
	default:
		panic(errors.New("::Marshal error: no matching type"))
	}

	return buf, err
}

// marshalString: marshal string to []byte
func marshalString(input interface{}) (buf []byte, err error) {
	return []byte(fmt.Sprintf("%v", input)), nil
}

// marshalInt32: marshal int32 to []byte
func marshalInt32(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarInt32(input.(int32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarInt32(buf, input.(int32))
	return buf, err
}

// marshalUint32: marshal uint32 to []byte
func marshalUint32(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarUInt32(input.(uint32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarUInt32(buf, input.(uint32))
	return buf, err
}

// marshalInt64: marshal int64 to []byte
func marshalInt64(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarInt64(input.(int64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarInt64(buf, input.(int64))
	return buf, err
}

// marshalUint64: marshal uint64 to []byte
func marshalUint64(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarUInt64(input.(uint64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarUInt64(buf, input.(uint64))
	return buf, err
}

// marshalFloat32: marshal float32 to []byte
func marshalFloat32(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfVarFloat32(input.(float32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodeVarFloat32(buf, input.(float32))
	return buf, err
}

// marshalFloat64: marshal float64 to []byte
func marshalFloat64(input interface{}) (buf []byte, err error) {
	size := encoding.SizeOfVarFloat64(input.(float64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodeVarFloat64(buf, input.(float64))
	return buf, err
}

// marshalStringSlice: marshal string slice to []byte
func marshalStringSlice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToStringSliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalInt32Slice: marshal int32 slice to []byte
func marshalInt32Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToInt64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalUint32Slice: marshal uint32 slice to []byte
func marshalUint32Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUInt64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt32Value(uint32(v.(uint64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalInt64Slice: marshal int64 slice to []byte
func marshalInt64Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToInt64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalUint64Slice: marshal uint64 slice to []byte
func marshalUint64Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUInt64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt64Value(v.(uint64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalFloat32Slice: marshal float32 slice to []byte
func marshalFloat32Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUFloat64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat32Value(float32(v.(float64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

// marshalFloat64Slice: marshal float64 slice to []byte
func marshalFloat64Slice(observe byte, input interface{}) []byte {
	var node = y3.NewNodeArrayPacketEncoder(int(observe))
	if out, ok := utils.ToUFloat64SliceArray(input); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat64Value(v.(float64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}
