package y3

import (
	"fmt"
	"reflect"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

type StructEncoder interface {
	Encode(input interface{}) ([]byte, error)
}

// structEncoder: for encode structure to packet
type structEncoder struct {
	config  *StructEncoderConfig
	observe byte
}

// StructEncoderConfig: config for Encoder
type StructEncoderConfig struct {
	ZeroFields bool
	Mold       interface{}
	TagName    string // Default value: yomo
}

func NewStructEncoder(observe byte, mold interface{}) StructEncoder {
	config := &StructEncoderConfig{
		ZeroFields: true,
		Mold:       mold,
	}

	return NewStructEncoderWithConfig(observe, config)
}

// NewStructEncoderWithConfig: create a Encoder for Struct
func NewStructEncoderWithConfig(observe byte, config *StructEncoderConfig) StructEncoder {
	if config.TagName == "" {
		config.TagName = "yomo"
	}

	result := &structEncoder{
		config:  config,
		observe: observe,
	}

	return result
}

func (e structEncoder) Encode(input interface{}) ([]byte, error) {
	packetEncoder, err := e.encode(input)
	if err != nil {
		return nil, err
	}

	if len(packetEncoder.GetValBuf()) == 0 {
		return nil, fmt.Errorf("::Encode valBuf of packetEncoder is empty")
	}

	buf := packetEncoder.Encode()
	return buf, err
}

// encode: encode interface to NodePacketEncoder
func (e *structEncoder) encode(input interface{}) (*NodePacketEncoder, error) {
	var inputVal reflect.Value

	if input != nil {
		inputVal = reflect.ValueOf(input)
		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if input == nil {
		return nil, fmt.Errorf("::encode input is nill")
	}

	if !inputVal.IsValid() {
		return nil, fmt.Errorf("::encode input value is not valid")
	}

	var packetEncode *NodePacketEncoder

	inputKind := inputVal.Kind()
	switch inputKind {
	case reflect.Struct:
		if utils.IsEmptyKey(e.observe) {
			packetEncode = e.encodeStruct(reflect.ValueOf(input), nil)
		} else {
			root := NewNodePacketEncoder(int(startingToken))
			packetEncode = e.encodeStruct(reflect.ValueOf(input), NewNodePacketEncoder(int(e.observe)))
			if !packetEncode.IsEmpty() {
				root.AddNodePacket(packetEncode)
			}
			packetEncode = root
		}
	case reflect.Slice:
		if utils.IsEmptyKey(e.observe) {
			packetEncode = e.encodeSlice(reflect.ValueOf(input), nil)
		} else {
			root := NewNodeArrayPacketEncoder(int(startingToken))
			packetEncode = e.encodeSlice(reflect.ValueOf(input), NewNodePacketEncoder(int(e.observe)))
			if !packetEncode.IsEmpty() {
				root.AddNodePacket(packetEncode)
			}
			packetEncode = root
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", inputKind)
	}

	return packetEncode, nil
}

// encodeSlice: encode slice to NodePacketEncoder
func (e *structEncoder) encodeSlice(sliceVal reflect.Value, root *NodePacketEncoder) *NodePacketEncoder {
	if root == nil {
		root = NewNodeArrayPacketEncoder(int(startingToken))
	}

	structType := sliceVal.Type()

	for i := 0; i < sliceVal.Len(); i++ {
		elemType := structType.Elem()
		switch elemType.Kind() {
		case reflect.Struct:
			currentValue := sliceVal.Index(i)
			p := e.encodeStruct(currentValue, NewNodePacketEncoder(utils.KeyOfArrayItem))
			if !p.IsEmpty() {
				root.AddNodePacket(p)
			}
		default:
			panic(fmt.Errorf("root slice unsupported type: %s", elemType.Kind()))
		}
	}

	return root
}

// encodeStruct: encode struct to NodePacketEncoder
func (e *structEncoder) encodeStruct(structVal reflect.Value, root *NodePacketEncoder) *NodePacketEncoder {
	if root == nil {
		root = NewNodePacketEncoder(int(startingToken))
	}

	var fields []field

	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fields = append(fields, field{fieldType, structVal.Field(i)})
	}

	for _, f := range fields {
		structField, fieldValue := f.field, f.val
		fieldName := fieldNameByTag(e.config.TagName, structField)
		e.encodeStructFromField(structField.Type, fieldName, fieldValue, root)
	}

	return root
}

// encodeStructFromField: encode struct from field
func (e *structEncoder) encodeStructFromField(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, en *NodePacketEncoder) {
	if fieldType.Kind() == reflect.Struct {
		leafNode := NewNodePacketEncoder(int(utils.KeyOf(fieldName)))
		fieldValueType := fieldValue.Type()
		for i := 0; i < fieldValueType.NumField(); i++ {
			thisFieldName := fieldNameByTag(e.config.TagName, fieldValueType.Field(i))
			thisFieldValue := fieldValue.Field(i)
			e.encodeStructFromField(thisFieldValue.Type(), thisFieldName, thisFieldValue, leafNode)
		}
		if !leafNode.IsEmpty() {
			en.AddNodePacket(leafNode)
		}
		return
	}

	var ppe = NewPrimitivePacketEncoder(int(utils.KeyOf(fieldName)))
	switch fieldType.Kind() {
	case reflect.String:
		ppe.SetStringValue(e.fieldValueToString(fieldType, fieldValue))
	case reflect.Int32:
		ppe.SetInt32Value(e.fieldValueToInt32(fieldType, fieldValue))
	case reflect.Uint32:
		ppe.SetUInt32Value(e.fieldValueToUint32(fieldType, fieldValue))
	case reflect.Int64:
		ppe.SetInt64Value(e.fieldValueToInt64(fieldType, fieldValue))
	case reflect.Uint64:
		ppe.SetUInt64Value(e.fieldValueToUInt64(fieldType, fieldValue))
	case reflect.Float32:
		ppe.SetFloat32Value(e.fieldValueToFloat32(fieldType, fieldValue))
	case reflect.Float64:
		ppe.SetFloat64Value(e.fieldValueToFloat64(fieldType, fieldValue))
	case reflect.Bool:
		ppe.SetBoolValue(e.fieldValueToBool(fieldType, fieldValue))
	case reflect.Array:
		arrNode := NewNodeArrayPacketEncoder(int(utils.KeyOf(fieldName)))
		e.encodeArrayFromField(fieldValue, arrNode)
		en.AddNodePacket(arrNode)
		return
	case reflect.Slice:
		sliceNode := NewNodeArrayPacketEncoder(int(utils.KeyOf(fieldName)))
		e.encodeSliceFromField(fieldValue, sliceNode)
		en.AddNodePacket(sliceNode)
		return
	default:
		panic(fmt.Errorf("there are no matches of any type: %v", fieldType.Kind()))
	}

	if !ppe.IsEmpty() {
		en.AddPrimitivePacket(ppe)
	}
}

// encodeArrayFromField: encode array from field
func (e *structEncoder) encodeArrayFromField(fieldValue reflect.Value, en *NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		e.encodeStructFromField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

// encodeSliceFromField: encode slice from field
func (e *structEncoder) encodeSliceFromField(fieldValue reflect.Value, en *NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		e.encodeStructFromField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

// fieldValueToString: get string value from fieldValue
func (e *structEncoder) fieldValueToString(fieldType reflect.Type, fieldValue reflect.Value) string {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).String()
	}
	return fieldValue.String()
}

// fieldValueToInt32: get int32 value from fieldValue
func (e *structEncoder) fieldValueToInt32(fieldType reflect.Type, fieldValue reflect.Value) int32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return int32(reflect.Zero(fieldType).Int())
	}
	return int32(fieldValue.Int())
}

// fieldValueToUint32: get uint32 value from fieldValue
func (e *structEncoder) fieldValueToUint32(fieldType reflect.Type, fieldValue reflect.Value) uint32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return uint32(reflect.Zero(fieldType).Uint())
	}
	return uint32(fieldValue.Uint())
}

// fieldValueToInt64: get int64 value from fieldValue
func (e *structEncoder) fieldValueToInt64(fieldType reflect.Type, fieldValue reflect.Value) int64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Int()
	}
	return fieldValue.Int()
}

// fieldValueToUInt64: get uint64 value from fieldValue
func (e *structEncoder) fieldValueToUInt64(fieldType reflect.Type, fieldValue reflect.Value) uint64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Uint()
	}
	return fieldValue.Uint()
}

// fieldValueToFloat32: get float32 value from fieldValue
func (e *structEncoder) fieldValueToFloat32(fieldType reflect.Type, fieldValue reflect.Value) float32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return float32(reflect.Zero(fieldType).Float())
	}
	return float32(fieldValue.Float())
}

// fieldValueToFloat64: get float64 value from fieldValue
func (e *structEncoder) fieldValueToFloat64(fieldType reflect.Type, fieldValue reflect.Value) float64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Float()
	}
	return fieldValue.Float()
}

// fieldValueToBool: get bool value from fieldValue
func (e *structEncoder) fieldValueToBool(fieldType reflect.Type, fieldValue reflect.Value) bool {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Bool()
	}
	return fieldValue.Bool()
}
