package y3

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// basicEncoder is a Encoder for BasicTestData data types
type basicEncoder interface {
	// Encode: encode interface to bytes
	Encode(input interface{}, signals ...*signal) (buf []byte, err error)
}

// basicEncoderImpl is implementation of the basicEncoder interface
type basicEncoderImpl struct {
	observe byte
	root    byte
}

// basicEncoderOption create basicEncoderImpl with option
type basicEncoderOption func(*basicEncoderImpl)

// basicEncoderOptionRoot set root value for creating basicEncoderImpl
func basicEncoderOptionRoot(root byte) basicEncoderOption {
	return func(b *basicEncoderImpl) {
		b.root = root
	}
}

// newBasicEncoder create a basicEncoder interface
func newBasicEncoder(observe byte, options ...func(*basicEncoderImpl)) basicEncoder {
	if utils.ForbiddenCustomizedKey(observe) {
		panic(fmt.Errorf("prohibit the use of this key: %#x", observe))
	}

	encoder := &basicEncoderImpl{observe: observe, root: utils.EmptyKey}

	for _, option := range options {
		option(encoder)
	}

	return encoder
}

// Encode encode interface{} to bytes
func (e *basicEncoderImpl) Encode(input interface{}, signals ...*signal) (buf []byte, err error) {
	encoders := make([]*PrimitivePacketEncoder, 0)
	for _, signal := range signals {
		encoders = append(encoders, signal.ToEncoder())
	}
	return e.encodeBasic(input, encoders)
}

// encodeBasic encode interface{} to bytes, and inserting signals
func (e *basicEncoderImpl) encodeBasic(input interface{}, signals []*PrimitivePacketEncoder) ([]byte, error) {
	if e.observe == 0 {
		panic(fmt.Errorf("observe cannot be 0"))
	}

	var primitiveEncoder *PrimitivePacketEncoder

	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.String:
		primitiveEncoder = e.encodeBasicString(input)
	case reflect.Int32:
		primitiveEncoder = e.encodeBasicInt32(input)
	case reflect.Uint32:
		primitiveEncoder = e.encodeBasicUint32(input)
	case reflect.Int64:
		primitiveEncoder = e.encodeBasicInt64(input)
	case reflect.Uint64:
		primitiveEncoder = e.encodeBasicUint64(input)
	case reflect.Float32:
		primitiveEncoder = e.encodeBasicFloat32(input)
	case reflect.Float64:
		primitiveEncoder = e.encodeBasicFloat64(input)
	case reflect.Bool:
		primitiveEncoder = e.encodeBasicBool(input)
	case reflect.Array, reflect.Slice:
		//e.marshalBasicSlice(value, e.root)
		return e.encodeBasicSlice(value, signals)
	default:
		panic(fmt.Errorf("marshal error, no matching type: %v", value.Kind()))
	}

	if primitiveEncoder == nil {
		panic("PrimitivePacketEncoder is nil")
	}

	if !utils.IsEmptyKey(e.root) {
		root := NewNodePacketEncoder(int(e.root))
		for _, signal := range signals {
			root.AddPrimitivePacket(signal)
		}
		root.AddPrimitivePacket(primitiveEncoder)
		return root.Encode(), nil
	}

	buf := make([][]byte, 0)
	for _, signal := range signals {
		buf = append(buf, signal.Encode())
	}
	buf = append(buf, primitiveEncoder.Encode())
	return bytes.Join(buf, []byte{}), nil
}

// encodeBasicSlice encode reflect.Value of slice, and inserting signals
func (e *basicEncoderImpl) encodeBasicSlice(value reflect.Value, signals []*PrimitivePacketEncoder) ([]byte, error) {
	if value.Len() == 0 {
		return nil, fmt.Errorf("no item is slice")
	}

	var nodeEncoder *NodePacketEncoder

	switch value.Index(0).Kind() {
	case reflect.String:
		nodeEncoder = e.encodeBasicStringSlice(value)
	case reflect.Int32:
		nodeEncoder = e.encodeBasicInt32Slice(value)
	case reflect.Uint32:
		nodeEncoder = e.encodeBasicUint32Slice(value)
	case reflect.Int64:
		nodeEncoder = e.encodeBasicInt64Slice(value)
	case reflect.Uint64:
		nodeEncoder = e.encodeBasicUint64Slice(value)
	case reflect.Float32:
		nodeEncoder = e.encodeBasicFloat32Slice(value)
	case reflect.Float64:
		nodeEncoder = e.encodeBasicFloat64Slice(value)
	case reflect.Bool:
		nodeEncoder = e.encodeBasicBoolSlice(value)
	default:
		panic(fmt.Errorf("marshal error, no matching type in SliceTestData: %v", value.Index(0).Kind()))
	}

	if nodeEncoder == nil {
		panic("NodePacketEncoder is nil")
	}

	if !utils.IsEmptyKey(e.root) {
		root := NewNodePacketEncoder(int(e.root))
		for _, signal := range signals {
			root.AddPrimitivePacket(signal)
		}
		root.AddNodePacket(nodeEncoder)
		return root.Encode(), nil
	}

	buf := make([][]byte, 0)
	for _, signal := range signals {
		buf = append(buf, signal.Encode())
	}
	buf = append(buf, nodeEncoder.Encode())
	return bytes.Join(buf, []byte{}), nil
}

// encodeBasicString encode string to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicString(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetStringValue(fmt.Sprintf("%v", input))
	return encoder
}

// encodeBasicInt32 encode int32 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicInt32(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetInt32Value(input.(int32))
	return encoder
}

// encodeBasicUint32 encode uint32 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicUint32(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetUInt32Value(input.(uint32))
	return encoder
}

// encodeBasicInt64 encode int64 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicInt64(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetInt64Value(input.(int64))
	return encoder
}

// encodeBasicUint64 encode uint64 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicUint64(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetUInt64Value(input.(uint64))
	return encoder
}

// encodeBasicFloat32 encode float32 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicFloat32(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetFloat32Value(input.(float32))
	return encoder
}

// encodeBasicFloat64 encode float64 to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicFloat64(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetFloat64Value(input.(float64))
	return encoder
}

// encodeBasicBool encode bool to PrimitivePacketEncoder
func (e *basicEncoderImpl) encodeBasicBool(input interface{}) *PrimitivePacketEncoder {
	var encoder = NewPrimitivePacketEncoder(int(e.observe))
	encoder.SetBoolValue(input.(bool))
	return encoder
}

// encodeBasicStringSlice encode reflect.Value of []string to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicStringSlice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToStringSliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicInt32Slice encode reflect.Value of []int32 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicInt32Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicUint32Slice encode reflect.Value of []uint32 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicUint32Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetUInt32Value(uint32(v.(uint64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicInt64Slice encode reflect.Value of []int64 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicInt64Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicUint64Slice encode reflect.Value of []uint64 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicUint64Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToUInt64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetUInt64Value(v.(uint64))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicFloat32Slice encode reflect.Value of []float32 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicFloat32Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetFloat32Value(float32(v.(float64)))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicFloat64Slice encode reflect.Value of []float64 to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicFloat64Slice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToUFloat64SliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetFloat64Value(v.(float64))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}

// encodeBasicBoolSlice encode reflect.Value of []bool to NodePacketEncoder
func (e *basicEncoderImpl) encodeBasicBoolSlice(value reflect.Value) *NodePacketEncoder {
	var node = NewNodeSlicePacketEncoder(int(e.observe))
	if out, ok := utils.ToBoolSliceArray(value.Interface()); ok {
		for _, v := range out {
			var item = NewPrimitivePacketEncoder(utils.KeyOfSliceItem)
			item.SetBoolValue(v.(bool))
			node.AddPrimitivePacket(item)
		}
	}
	return node
}
