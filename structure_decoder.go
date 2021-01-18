package y3

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

type StructDecoder interface {
	Decode(input []byte) (interface{}, error)
}

// structDecoder: for decode packet to structure
type structDecoder struct {
	config *StructDecoderConfig
}

// StructDecoderConfig: config for structDecoder
type StructDecoderConfig struct {
	ZeroFields bool // 在解码值前进行清零或清空操作
	Result     interface{}
	TagName    string // 默认值: yomo
}

// NewStructDecoder: create a structDecoder
func NewStructDecoder(mold interface{}) StructDecoder {
	config := &StructDecoderConfig{
		ZeroFields: true,
		Result:     mold,
	}

	decoder, err := NewStructDecoderWithConfig(config)
	if err != nil {
		panic(fmt.Errorf("NewStructDecoderWithConfig error: %v", err))
	}

	return decoder
}

// NewStructDecoderWithConfig: create a structDecoder with StructDecoderConfig
func NewStructDecoderWithConfig(config *StructDecoderConfig) (StructDecoder, error) {
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return nil, errors.New("result must be addressable (a pointer)")
	}

	if config.TagName == "" {
		config.TagName = "yomo"
	}

	result := &structDecoder{
		config: config,
	}

	return result, nil
}

func (d *structDecoder) Decode(input []byte) (interface{}, error) {
	err := d.decode(input)
	if err != nil {
		return nil, err
	}
	return d.config.Result, nil
}

// Decode: public func for decode
func (d *structDecoder) decode(input []byte) error {
	node, _, err := DecodeNodePacket(input)
	if err != nil {
		return err
	}

	return d.decodeNode(node, reflect.ValueOf(d.config.Result).Elem())
}

// decodeNode: observe NodePacket to struct
func (d *structDecoder) decodeNode(input *NodePacket, outVal reflect.Value) error {
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

// decodeSlice: decode slice
func (d *structDecoder) decodeSlice(data *NodePacket, val reflect.Value) error {
	elemType := val.Type().Elem()

	items := make([]reflect.Value, 0)
	for _, node := range data.NodePackets {
		elemValue := reflect.Indirect(reflect.New(elemType))
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

// decodeStruct: decode struct
func (d *structDecoder) decodeStruct(data *NodePacket, val reflect.Value) error {
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

// decodeStructFromNodePacket: decode struct from NodePacket
func (d *structDecoder) decodeStructFromNodePacket(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, dataVal *NodePacket) error {
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

	key := utils.KeyOf(fieldName)
	if !d.allowCustomizedKey(key) && key != startingToken {
		return fmt.Errorf("not allow key:%#x", key)
	}

	return nil
}

// getKind: get kind of Value
func (d *structDecoder) getKind(val reflect.Value) reflect.Kind {
	return val.Kind()
}

// takeValueByKey: take Value by fieldName
func (d *structDecoder) takeValueByKey(fieldName string, fieldType reflect.Type, node *NodePacket) (reflect.Value, bool) {
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

// takeNodeValue: take Value from NodePacket
func (d *structDecoder) takeNodeValue(fieldType reflect.Type, nodePacket NodePacket) (reflect.Value, bool) {
	if nodePacket.IsArray() {
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

// paddingToStruct: padding to struct from NodePacket
func (d *structDecoder) paddingToStruct(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
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

// paddingToArray: padding to Array from NodePacket
func (d *structDecoder) paddingToArray(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
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
	reflect.Copy(arrayValue, slice) // TODO: 注意性能影响
	return arrayValue
}

// paddingToSlice: padding to Slice from NodePacket
func (d *structDecoder) paddingToSlice(fieldType reflect.Type, nodePacket NodePacket) reflect.Value {
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

// takePrimitiveValue: take primitive value from PrimitivePacket
func (d *structDecoder) takePrimitiveValue(fieldType reflect.Type, primitivePacket PrimitivePacket) (reflect.Value, bool) {
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

// matchingKey: matching node of key from NodePacket
func (d *structDecoder) matchingKey(key byte, node *NodePacket) (flag bool, isNode bool, packet interface{}) {
	if len(node.PrimitivePackets) > 0 {
		for _, p := range node.PrimitivePackets {
			if key == p.SeqID() {
				return true, false, p
			}
		}
	}

	if len(node.NodePackets) > 0 {
		for _, n := range node.NodePackets {
			if key == n.SeqID() {
				return true, true, n
			}
			//return matchingKey(key, &n)
			flag, isNode, packet = d.matchingKey(key, &n)
			if flag {
				return
			}
		}
	}

	return false, false, nil
}

// allowCustomizedKey: allow customized key
func (d *structDecoder) allowCustomizedKey(key byte) bool {
	switch key {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return false
	}
	return true
}
