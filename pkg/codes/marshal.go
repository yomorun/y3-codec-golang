package codes

import (
	"errors"
	"fmt"
	"reflect"

	y3 "github.com/yomorun/yomo-codec-golang"
	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// interface to []byte, serialization
func (codec *yomoCodec) Marshal(T interface{}) (buf []byte, err error) {
	//fmt.Printf("#74 reflect.ValueOf(T)=%v\n", reflect.ValueOf(T))
	switch reflect.ValueOf(T).Kind() {
	case reflect.String:
		buf, err = marshalString(T)
	case reflect.Int32:
		buf, err = marshalInt32(T)
	case reflect.Uint32:
		buf, err = marshalUint32(T)
	case reflect.Int64:
		buf, err = marshalInt64(T)
	case reflect.Uint64:
		buf, err = marshalUint64(T)
	case reflect.Float32:
		buf, err = marshalFloat32(T)
	case reflect.Float64:
		buf, err = marshalFloat64(T)
	case reflect.Array, reflect.Slice:
		if reflect.ValueOf(T).Len() == 0 {
			break
		}
		switch reflect.ValueOf(T).Index(0).Elem().Kind() {
		case reflect.String:
			buf = marshalStringSlice(codec.Observe, T)
		case reflect.Int32:
			buf = marshalInt32Slice(codec.Observe, T)
		case reflect.Uint32:
			buf = marshalUint32Slice(codec.Observe, T)
		case reflect.Int64:
			buf = marshalInt64Slice(codec.Observe, T)
		case reflect.Uint64:
			buf = marshalUint64Slice(codec.Observe, T)
		case reflect.Float32:
			buf = marshalFloat32Slice(codec.Observe, T)
		case reflect.Float64:
			buf = marshalFloat64Slice(codec.Observe, T)
		default:
			panic(errors.New("::Marshal error: no matching type in Slice"))
		}
	default:
		panic(errors.New("::Marshal error: no matching type"))
	}

	return buf, err
}

func marshalString(T interface{}) (buf []byte, err error) {
	return []byte(fmt.Sprintf("%v", T)), nil
}

func marshalInt32(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarInt32(T.(int32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarInt32(buf, T.(int32))
	return buf, err
}

func marshalUint32(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarUInt32(T.(uint32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarUInt32(buf, T.(uint32))
	return buf, err
}

func marshalInt64(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarInt64(T.(int64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarInt64(buf, T.(int64))
	return buf, err
}

func marshalUint64(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfPVarUInt64(T.(uint64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodePVarUInt64(buf, T.(uint64))
	return buf, err
}

func marshalFloat32(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfVarFloat32(T.(float32))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodeVarFloat32(buf, T.(float32))
	return buf, err
}

func marshalFloat64(T interface{}) (buf []byte, err error) {
	size := encoding.SizeOfVarFloat64(T.(float64))
	codec := encoding.VarCodec{Size: size}
	buf = make([]byte, size)
	err = codec.EncodeVarFloat64(buf, T.(float64))
	return buf, err
}

func marshalStringSlice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToStringSliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalInt32Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToInt64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalUint32Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToUInt64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt32Value(uint32(v.(uint64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalInt64Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToInt64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalUint64Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToUInt64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetUInt64Value(v.(uint64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalFloat32Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToUFloat64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat32Value(float32(v.(float64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}

func marshalFloat64Slice(observe string, T interface{}) []byte {
	key := keyOf(observe)
	var node = y3.NewNodeArrayPacketEncoder(int(key))
	if out, ok := utils.ToUFloat64SliceArray(T); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(utils.KeyOfArrayItem)
			item.SetFloat64Value(v.(float64))
			node.AddPrimitivePacket(item)
		}
	}
	return node.GetValBuf()
}
