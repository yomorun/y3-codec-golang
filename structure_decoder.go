package y3

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// structEncoder is a Decoder for Struct type
type structDecoder interface {
	// Decode decode bytes to interface
	Decode(input []byte) (interface{}, error)
}

// structDecoderImpl is implementation of the structDecoder interface
type structDecoderImpl struct {
	config *structDecoderConfig
	result interface{}
}

// structDecoderConfig is configuration for structDecoderImpl
type structDecoderConfig struct {
	ZeroFields bool   // 在解码值前进行清零或清空操作
	TagName    string // 默认值: yomo
}

// structDecoderOption create structDecoderImpl with option
type structDecoderOption func(*structDecoderImpl)

// structDecoderOptionConfig set structDecoderConfig value for creating structDecoderImpl
func structDecoderOptionConfig(config *structDecoderConfig) structDecoderOption {
	return func(e *structDecoderImpl) {
		e.config = config
	}
}

// newStructDecoder create a structDecoder interface
func newStructDecoder(mold interface{}, options ...func(*structDecoderImpl)) structDecoder {
	// check mold
	val := reflect.ValueOf(mold)
	if val.Kind() != reflect.Ptr {
		panic(errors.New("mold must be a pointer"))
	}
	val = val.Elem()
	if !val.CanAddr() {
		panic(errors.New("mold must be addressable (a pointer)"))
	}

	decoder := &structDecoderImpl{
		config: &structDecoderConfig{
			ZeroFields: true,
			TagName:    "y3",
		},
		result: mold,
	}

	for _, option := range options {
		option(decoder)
	}

	return decoder
}

// Decode decode bytes to interface
func (d *structDecoderImpl) Decode(input []byte) (interface{}, error) {
	err := d.decode(input)
	if err != nil {
		return nil, err
	}
	return d.result, nil
}

// Decode decode interface to d.config.Result
func (d *structDecoderImpl) decode(input []byte) error {
	node, _, err := DecodeNodePacket(input)
	if err != nil {
		return err
	}

	return d.decodeNode(node, reflect.ValueOf(d.result).Elem())
}

// decodeNode is entry function for decoding NodePacket
func (d *structDecoderImpl) decodeNode(input *NodePacket, outVal reflect.Value) error {
	var inputVal reflect.Value

	if input != nil {
		inputVal = reflect.ValueOf(input)

		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if input == nil {
		// 如果数据为nil，那么不设置任何值，除非ZeroField为true
		if d.config.ZeroFields {
			outVal.Set(reflect.Zero(outVal.Type()))
		}
		return nil
	}

	if !inputVal.IsValid() {
		// 如果数据为invalid，则仅设置期zero值
		outVal.Set(reflect.Zero(outVal.Type()))
		return nil
	}

	var err error
	outputKind := d.getKind(outVal)
	switch outputKind {
	case reflect.Struct:
		err = d.decodeStruct(input, outVal)
	case reflect.Slice:
		err = d.decodeSlice(input, outVal)
	default:
		return fmt.Errorf("decode unsupported type: %s", outputKind)
	}

	return err
}

// decodeSlice decode slice of struct
func (d *structDecoderImpl) decodeSlice(data *NodePacket, val reflect.Value) error {
	elemType := val.Type().Elem()

	items := make([]reflect.Value, 0)
	for i := range data.NodePackets {
		elemValue := reflect.Indirect(reflect.New(elemType))
		node := data.NodePackets[i]
		err := d.decodeStruct(&node, elemValue)
		if err != nil {
			return err
		}
		items = append(items, elemValue)
	}

	sliceType := reflect.SliceOf(elemType)
	sliceValue := reflect.New(sliceType).Elem()
	sliceValue = reflect.Append(sliceValue, items...)
	val.Set(sliceValue)

	if len(data.PrimitivePackets) > 0 {
		panic(fmt.Errorf("root primitive slice is unsupported: PrimitivePackets-len=%v", len(data.PrimitivePackets)))
	}

	return nil
}

// decodeStruct decode struct
func (d *structDecoderImpl) decodeStruct(data *NodePacket, val reflect.Value) error {
	var fields []field
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		structField := valType.Field(i)
		fields = append(fields, field{structField, val.Field(i)})
	}

	for _, f := range fields {
		structField, fieldValue := f.field, f.val
		err := d.decodeStructFromNodePacket(structField.Type, fieldNameByTag(d.config.TagName, structField), fieldValue, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// decodeStructFromNodePacket decode struct from NodePacket
func (d *structDecoderImpl) decodeStructFromNodePacket(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, dataVal *NodePacket) error {
	if fieldType.Kind() == reflect.Struct {
		fieldValueType := fieldValue.Type()
		for i := 0; i < fieldValueType.NumField(); i++ {
			currentFieldName := fieldNameByTag(d.config.TagName, fieldValueType.Field(i))
			currentFieldValue := fieldValue.Field(i)
			err := d.decodeStructFromNodePacket(currentFieldValue.Type(), currentFieldName, currentFieldValue, dataVal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	obtainedValue, flag := d.takeValueByKey(fieldName, fieldType, dataVal)
	if !flag {
		if d.config.ZeroFields == false {
			return fmt.Errorf("not fond fieldName:%#x", fieldName)
		}
		// 用空值填充找不到的字段
		obtainedValue = reflect.Zero(fieldType)
	}

	if !fieldValue.IsValid() {
		// This should never happen
		panic("structField is not valid")
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("fieldValue can not set:%v", fieldValue)
	}

	fieldValue.Set(obtainedValue)

	return nil
}

// getKind get kind of reflect.Value
func (d *structDecoderImpl) getKind(val reflect.Value) reflect.Kind {
	return val.Kind()
}

// takeValueByKey take Value by fieldName
func (d *structDecoderImpl) takeValueByKey(fieldName string, fieldType reflect.Type, node *NodePacket) (reflect.Value, bool) {
	key := utils.KeyOf(fieldName)
	flag, isNode, packet := d.matchingKey(key, node)
	if flag == false {
		return reflect.Indirect(reflect.ValueOf(packet)), false
	}
	if isNode {
		nodePacket := packet.(NodePacket)
		return d.takeNodeValue(fieldType, nodePacket)
	}

	primitivePacket := packet.(PrimitivePacket)
	return d.takePrimitiveValue(fieldType, primitivePacket)
}

// takeNodeValue take Value from NodePacket
func (d *structDecoderImpl) takeNodeValue(fieldType reflect.Type, nodePacket NodePacket) (reflect.Value, bool) {
	if nodePacket.IsSlice() {
		switch fieldType.Kind() {
		case reflect.Array:
			return d.paddingToArray(fieldType, nodePacket), true
		case reflect.Slice:
			return d.paddingToSlice(fieldType, nodePacket), true
		default:
			panic(fmt.Errorf("unimplemented type %v", fieldType.Kind()))
		}
	}

	if len(nodePacket.PrimitivePackets) > 0 || len(nodePacket.NodePackets) > 0 {
		switch fieldType.Kind() {
		case reflect.Struct:
			return d.paddingToStruct(fieldType, nodePacket), true
		default:
			panic(fmt.Errorf("unimplemented type %v when have PrimitivePackets", fieldType.Kind()))
		}
	}

	return reflect.Indirect(reflect.ValueOf(nodePacket)), true
}

// paddingToStruct padding to struct from NodePacket
func (d *structDecoderImpl) paddingToStruct(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
	obj := reflect.New(fieldType)
	obj = reflect.Indirect(obj)

	for _, v := range nodePacket.PrimitivePackets {
		for i := 0; i < obj.NumField(); i++ {
			structField := obj.Type().Field(i)
			valueField := obj.Field(i)
			fieldName := fieldNameByTag(d.config.TagName, structField)
			if v.SeqID() == utils.KeyOf(fieldName) {
				vv, _ := d.takePrimitiveValue(valueField.Type(), v)
				valueField.Set(vv)
			}
		}
	}

	for _, v := range nodePacket.NodePackets {
		for i := 0; i < obj.NumField(); i++ {
			structField := obj.Type().Field(i)
			valueField := obj.Field(i)
			fieldName := fieldNameByTag(d.config.TagName, structField)
			if v.SeqID() == utils.KeyOf(fieldName) {
				vv, _ := d.takeNodeValue(structField.Type, v)
				valueField.Set(vv)
			}
		}
	}

	return reflect.Indirect(obj)
}

// paddingToArray padding to ArrayTestData from NodePacket
func (d *structDecoderImpl) paddingToArray(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
	sliceType := reflect.SliceOf(fieldType.Elem())
	sliceValue := reflect.New(sliceType).Elem()

	if len(nodePacket.NodePackets) > 0 {
		items := make([]reflect.Value, 0)
		for _, node := range nodePacket.NodePackets {
			itemValue, _ := d.takeNodeValue(fieldType.Elem(), node)
			items = append(items, itemValue)
		}
		slice := reflect.Append(sliceValue, items...)

		arrayType := reflect.ArrayOf(len(nodePacket.NodePackets), fieldType.Elem())
		arrayValue := reflect.New(arrayType).Elem()
		reflect.Copy(arrayValue, slice)
		return arrayValue
	}

	items := make([]reflect.Value, 0)
	for _, primitivePacket := range nodePacket.PrimitivePackets {
		itemValue, _ := d.takePrimitiveValue(fieldType.Elem(), primitivePacket)
		items = append(items, itemValue)
	}
	slice := reflect.Append(sliceValue, items...)

	arrayType := reflect.ArrayOf(len(nodePacket.PrimitivePackets), fieldType.Elem())
	arrayValue := reflect.New(arrayType).Elem()
	reflect.Copy(arrayValue, slice)
	return arrayValue
}

// paddingToSlice padding to SliceTestData from NodePacket
func (d *structDecoderImpl) paddingToSlice(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
	sliceType := reflect.SliceOf(fieldType.Elem())
	sliceValue := reflect.New(sliceType).Elem()

	if len(nodePacket.NodePackets) > 0 {
		items := make([]reflect.Value, 0)
		for _, node := range nodePacket.NodePackets {
			itemValue, _ := d.takeNodeValue(fieldType.Elem(), node)
			items = append(items, itemValue)
		}
		return reflect.Append(sliceValue, items...)
	}

	items := make([]reflect.Value, 0)
	for _, primitivePacket := range nodePacket.PrimitivePackets {
		itemValue, _ := d.takePrimitiveValue(fieldType.Elem(), primitivePacket)
		items = append(items, itemValue)
	}
	slice := reflect.Append(sliceValue, items...)

	return slice
}

// takePrimitiveValue take primitive value from PrimitivePacket
func (d *structDecoderImpl) takePrimitiveValue(fieldType reflect.Type, primitivePacket PrimitivePacket) (reflect.Value, bool) {
	switch fieldType.Kind() {
	case reflect.String:
		val, err := primitivePacket.ToUTF8String()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Int32:
		val, err := primitivePacket.ToInt32()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Int64:
		val, err := primitivePacket.ToInt64()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Uint32:
		val, err := primitivePacket.ToUInt32()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Uint64:
		val, err := primitivePacket.ToUInt64()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Float32:
		val, err := primitivePacket.ToFloat32()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Float64:
		val, err := primitivePacket.ToFloat64()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	case reflect.Bool:
		val, err := primitivePacket.ToBool()
		if err != nil {
			panic(err)
		}
		return reflect.Indirect(reflect.ValueOf(val)), true
	default:
		panic(errors.New("::takeValueByKey error: no matching type"))
	}
}

// matchingKey matching node of key from NodePacket
func (d *structDecoderImpl) matchingKey(key byte, node *NodePacket) (flag bool, isNode bool, packet interface{}) {
	if len(node.PrimitivePackets) > 0 {
		for _, p := range node.PrimitivePackets {
			if key == p.SeqID() {
				return true, false, p
			}
		}
	}

	if len(node.NodePackets) > 0 {
		for i := range node.NodePackets {
			n := node.NodePackets[i]
			if key == n.SeqID() {
				return true, true, n
			}
			flag, isNode, packet = d.matchingKey(key, &n)
			if flag {
				return
			}
		}
	}

	return false, false, nil
}

var (
	// rootToken: mark the root node
	rootToken byte = 0x01
)

// field store the contents of a reflect
type field struct {
	field reflect.StructField
	val   reflect.Value
}

// fieldNameByTag: get fieldName by considering tagName
func fieldNameByTag(tagName string, field reflect.StructField) string {
	fieldName := field.Name

	tagValue := field.Tag.Get(tagName)
	tagValue = strings.SplitN(tagValue, ",", 2)[0]
	if tagValue != "" {
		fieldName = tagValue
	}

	return fieldName
}
