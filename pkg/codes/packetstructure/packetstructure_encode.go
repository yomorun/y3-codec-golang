package packetstructure

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	y3 "github.com/yomorun/yomo-codec-golang"
)

type Encoder struct {
	config *EncoderConfig
}

type EncoderConfig struct {
	Mold    interface{}
	TagName string // 默认值: yomo
}

func NewEncoder(config *EncoderConfig) (*Encoder, error) {
	if config.TagName == "" {
		config.TagName = "yomo"
	}

	result := &Encoder{
		config: config,
	}

	return result, nil
}

func Encode(input interface{}) (*y3.NodePacket, error) {
	config := &EncoderConfig{
		Mold: input,
	}

	encoder, err := NewEncoder(config)
	if err != nil {
		return nil, err
	}

	return encoder.Encode(input)
}

func (d *Encoder) Encode(input interface{}) (*y3.NodePacket, error) {
	packetEncoder, err := d.encode(input)
	if err != nil {
		return nil, err
	}

	if len(packetEncoder.GetValBuf()) == 0 {
		return nil, fmt.Errorf("::Encode valBuf of packetEncoder is empty")
	}

	buf := packetEncoder.Encode()
	node, _, err := y3.DecodeNodePacket(buf)
	return node, err
}

func (d *Encoder) encode(input interface{}) (*y3.NodePacketEncoder, error) {
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
		packetEncode = d.encodeStruct(input, nil)
	default:
		return nil, fmt.Errorf("unsupported type: %s", inputKind)
	}

	return packetEncode, nil
}

func (d *Encoder) encodeStruct(input interface{}, root *y3.NodePacketEncoder) *y3.NodePacketEncoder {
	if root == nil {
		root = y3.NewNodePacketEncoder(int(startingToken))
	}

	var fields []field

	structVal := reflect.ValueOf(input)
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fields = append(fields, field{fieldType, structVal.Field(i)})
	}

	for _, f := range fields {
		structField, fieldValue := f.field, f.val
		fieldName := d.fieldNameByTag(structField)
		d.encodeStructFromStructField(structField.Type, fieldName, fieldValue, root)
	}

	return root
}

func (d *Encoder) encodeStructFromStructField(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	if fieldType.Kind() == reflect.Struct {
		leafNode := y3.NewNodePacketEncoder(int(keyOf(fieldName)))
		fieldValueType := fieldValue.Type()
		for i := 0; i < fieldValueType.NumField(); i++ {
			thisFieldName := d.fieldNameByTag(fieldValueType.Field(i))
			thisFieldValue := fieldValue.Field(i)
			d.encodeStructFromStructField(thisFieldValue.Type(), thisFieldName, thisFieldValue, leafNode)
		}
		en.AddNodePacket(leafNode)
		return
	}

	var ppe = y3.NewPrimitivePacketEncoder(int(keyOf(fieldName)))
	switch fieldType.Kind() {
	case reflect.String:
		ppe.SetStringValue(fieldValue.String())
	case reflect.Int32:
		ppe.SetInt32Value(int32(fieldValue.Int()))
	case reflect.Uint32:
		ppe.SetUInt32Value(uint32(fieldValue.Uint()))
	case reflect.Int64:
		ppe.SetInt64Value(fieldValue.Int())
	case reflect.Uint64:
		ppe.SetUInt64Value(fieldValue.Uint())
	case reflect.Float32:
		ppe.SetFloat32Value(float32(fieldValue.Float()))
	case reflect.Float64:
		ppe.SetFloat64Value(fieldValue.Float())
	case reflect.Array:
		arrNode := y3.NewNodeArrayPacketEncoder(int(keyOf(fieldName)))
		d.decodeArray(fieldValue, arrNode)
		en.AddNodePacket(arrNode)
		return
	case reflect.Slice:
		sliceNode := y3.NewNodeArrayPacketEncoder(int(keyOf(fieldName)))
		d.decodeArray(fieldValue, sliceNode)
		en.AddNodePacket(sliceNode)
		return
	}

	//fmt.Printf("#306 kind=%v\n", fieldType.Kind())
	en.AddPrimitivePacket(ppe)
}

func (d *Encoder) decodeArray(fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		d.encodeStructFromStructField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

func (d *Encoder) decodeSlice(fieldValue reflect.Value, en *y3.NodePacketEncoder) {
	for i := 0; i < fieldValue.Len(); i++ {
		currentData := fieldValue.Index(i)
		currentType := fieldValue.Index(i).Type()
		d.encodeStructFromStructField(currentType, utils.KeyStringOfArrayItem, currentData, en)
	}
}

func (d *Encoder) fieldNameByTag(field reflect.StructField) string {
	fieldName := field.Name

	tagValue := field.Tag.Get(d.config.TagName)
	tagValue = strings.SplitN(tagValue, ",", 2)[0]
	if tagValue != "" {
		fieldName = tagValue
	}

	return fieldName
}
