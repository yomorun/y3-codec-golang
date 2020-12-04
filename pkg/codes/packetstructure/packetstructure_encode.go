package packetstructure

import (
	"fmt"
	"reflect"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	y3 "github.com/yomorun/yomo-codec-golang"
)

// Encoder: for encode structure to packet
type Encoder struct {
	config *EncoderConfig
}

// EncoderConfig: config for Encoder
type EncoderConfig struct {
	ZeroFields bool
	Mold       interface{}
	TagName    string // 默认值: yomo
}

// NewEncoder: create a Encoder
func NewEncoder(config *EncoderConfig) (*Encoder, error) {
	if config.TagName == "" {
		config.ZeroFields = true
		config.TagName = "yomo"
	}

	result := &Encoder{
		config: config,
	}

	return result, nil
}

// defaultEncoder: create a default Encoder
func defaultEncoder(input interface{}) (*Encoder, error) {
	config := &EncoderConfig{
		Mold: input,
	}

	return NewEncoder(config)
}

// Encode: shortcut of Encoder, return NodePacket
func Encode(input interface{}) (*y3.NodePacket, error) {
	return EncodeWith(packetutils.KeyOf(""), input)
}

// EncodeWith: shortcut of Encoder, with observe, return NodePacket
func EncodeWith(observe byte, input interface{}) (*y3.NodePacket, error) {
	encoder, err := defaultEncoder(input)
	if err != nil {
		return nil, err
	}

	return encoder.Encode(observe, input)
}

// EncodeToBytesWith: shortcut of Encoder, return []byte
func EncodeToBytesWith(observe byte, input interface{}) ([]byte, error) {
	encoder, err := defaultEncoder(input)
	if err != nil {
		return nil, err
	}
	return encoder.EncodeToBytes(observe, input)
}

// Encode: public func for encode, return NodePacket
func (e *Encoder) Encode(observe byte, input interface{}) (*y3.NodePacket, error) {
	buf, err := e.EncodeToBytes(observe, input)
	if err != nil {
		return nil, err
	}

	node, _, err := y3.DecodeNodePacket(buf)
	return node, err
}

// EncodeToBytes: encode to bytes
func (e *Encoder) EncodeToBytes(observe byte, input interface{}) ([]byte, error) {
	packetEncoder, err := e.encode(observe, input)
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
func (e *Encoder) encode(observe byte, input interface{}) (*y3.NodePacketEncoder, error) {
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

	var packetEncode *y3.NodePacketEncoder

	inputKind := inputVal.Kind()
	switch inputKind {
	case reflect.Struct:
		if packetutils.IsEmptyKey(observe) {
			packetEncode = e.encodeStruct(reflect.ValueOf(input), nil)
		} else {
			root := y3.NewNodePacketEncoder(int(startingToken))
			packetEncode = e.encodeStruct(reflect.ValueOf(input), y3.NewNodePacketEncoder(int(observe)))
			root.AddNodePacket(packetEncode)
			packetEncode = root
		}
	case reflect.Slice:
		if packetutils.IsEmptyKey(observe) {
			packetEncode = e.encodeSlice(reflect.ValueOf(input), nil)
		} else {
			root := y3.NewNodeArrayPacketEncoder(int(startingToken))
			packetEncode = e.encodeSlice(reflect.ValueOf(input), y3.NewNodePacketEncoder(int(observe)))
			root.AddNodePacket(packetEncode)
			packetEncode = root
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", inputKind)
	}

	return packetEncode, nil
}

// encodeSlice: encode slice to NodePacketEncoder
func (e *Encoder) encodeSlice(sliceVal reflect.Value, root *y3.NodePacketEncoder) *y3.NodePacketEncoder {
	if root == nil {
		root = y3.NewNodeArrayPacketEncoder(int(startingToken))
	}

	structType := sliceVal.Type()

	for i := 0; i < sliceVal.Len(); i++ {
		elemType := structType.Elem()
		switch elemType.Kind() {
		case reflect.Struct:
			currentValue := sliceVal.Index(i)
			p := e.encodeStruct(currentValue, y3.NewNodePacketEncoder(utils.KeyOfArrayItem))
			root.AddNodePacket(p)
		default:
			panic(fmt.Errorf("root slice unsupported type: %s", elemType.Kind()))
		}
	}

	return root
}

// encodeStruct: encode struct to NodePacketEncoder
func (e *Encoder) encodeStruct(structVal reflect.Value, root *y3.NodePacketEncoder) *y3.NodePacketEncoder {
	if root == nil {
		root = y3.NewNodePacketEncoder(int(startingToken))
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
func (e *Encoder) encodeStructFromField(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	if fieldType.Kind() == reflect.Struct {
		leafNode := y3.NewNodePacketEncoder(int(packetutils.KeyOf(fieldName)))
		fieldValueType := fieldValue.Type()
		for i := 0; i < fieldValueType.NumField(); i++ {
			thisFieldName := fieldNameByTag(e.config.TagName, fieldValueType.Field(i))
			thisFieldValue := fieldValue.Field(i)
			e.encodeStructFromField(thisFieldValue.Type(), thisFieldName, thisFieldValue, leafNode)
		}
		en.AddNodePacket(leafNode)
		return
	}

	var ppe = y3.NewPrimitivePacketEncoder(int(packetutils.KeyOf(fieldName)))
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
		arrNode := y3.NewNodeArrayPacketEncoder(int(packetutils.KeyOf(fieldName)))
		e.encodeArrayFromField(fieldValue, arrNode)
		en.AddNodePacket(arrNode)
		return
	case reflect.Slice:
		sliceNode := y3.NewNodeArrayPacketEncoder(int(packetutils.KeyOf(fieldName)))
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
func (e *Encoder) encodeArrayFromField(fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		e.encodeStructFromField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

// encodeSliceFromField: encode slice from field
func (e *Encoder) encodeSliceFromField(fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		e.encodeStructFromField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

// fieldValueToString: get string value from fieldValue
func (e *Encoder) fieldValueToString(fieldType reflect.Type, fieldValue reflect.Value) string {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).String()
	}
	return fieldValue.String()
}

// fieldValueToInt32: get int32 value from fieldValue
func (e *Encoder) fieldValueToInt32(fieldType reflect.Type, fieldValue reflect.Value) int32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return int32(reflect.Zero(fieldType).Int())
	}
	return int32(fieldValue.Int())
}

// fieldValueToUint32: get uint32 value from fieldValue
func (e *Encoder) fieldValueToUint32(fieldType reflect.Type, fieldValue reflect.Value) uint32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return uint32(reflect.Zero(fieldType).Uint())
	}
	return uint32(fieldValue.Uint())
}

// fieldValueToInt64: get int64 value from fieldValue
func (e *Encoder) fieldValueToInt64(fieldType reflect.Type, fieldValue reflect.Value) int64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Int()
	}
	return fieldValue.Int()
}

// fieldValueToUInt64: get uint64 value from fieldValue
func (e *Encoder) fieldValueToUInt64(fieldType reflect.Type, fieldValue reflect.Value) uint64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Uint()
	}
	return fieldValue.Uint()
}

// fieldValueToFloat32: get float32 value from fieldValue
func (e *Encoder) fieldValueToFloat32(fieldType reflect.Type, fieldValue reflect.Value) float32 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return float32(reflect.Zero(fieldType).Float())
	}
	return float32(fieldValue.Float())
}

// fieldValueToFloat64: get float64 value from fieldValue
func (e *Encoder) fieldValueToFloat64(fieldType reflect.Type, fieldValue reflect.Value) float64 {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Float()
	}
	return fieldValue.Float()
}

// fieldValueToBool: get bool value from fieldValue
func (e *Encoder) fieldValueToBool(fieldType reflect.Type, fieldValue reflect.Value) bool {
	if fieldValue.IsZero() && e.config.ZeroFields {
		return reflect.Zero(fieldType).Bool()
	}
	return fieldValue.Bool()
}
