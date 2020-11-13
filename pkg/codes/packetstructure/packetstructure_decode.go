package packetstructure

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	y3 "github.com/yomorun/yomo-codec-golang"
)

type Decoder struct {
	config *DecoderConfig
}

type DecoderConfig struct {
	ZeroFields bool // 在解码值前进行清零或清空操作
	Result     interface{}
	TagName    string // 默认值: yomo
}

func NewDecoder(config *DecoderConfig) (*Decoder, error) {
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

	result := &Decoder{
		config: config,
	}

	return result, nil
}

func Decode(input *y3.NodePacket, output interface{}) error {
	config := &DecoderConfig{
		Result: output,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func (d *Decoder) Decode(input *y3.NodePacket) error {
	return d.decode(startingToken, input, reflect.ValueOf(d.config.Result).Elem())
}

// decode name: observe key, input: input packet, outVal:
func (d *Decoder) decode(name byte, input *y3.NodePacket, outVal reflect.Value) error {
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
	//case reflect.String:
	//	err = d.decodeString(name, input, outVal)
	case reflect.Struct:
		err = d.decodeStruct(name, input, outVal)
	default:
		return fmt.Errorf("%#x: unsupported type: %s", name, outputKind)
	}

	return err
}

func (d *Decoder) decodeStruct(name byte, data *y3.NodePacket, val reflect.Value) error {
	var fields []field
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		structField := valType.Field(i)
		fields = append(fields, field{structField, val.Field(i)})
	}
	//fmt.Printf("#202 structs len=%v\n", len(fields))

	for _, f := range fields {
		structField, fieldValue := f.field, f.val
		err := d.decodeStructFromNodePacket(structField.Type, d.fieldNameByTag(structField), fieldValue, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeStructFromNodePacket(fieldType reflect.Type, fieldName string, fieldValue reflect.Value, dataVal *y3.NodePacket) error {
	if fieldType.Kind() == reflect.Struct {
		fieldValueType := fieldValue.Type()
		for i := 0; i < fieldValueType.NumField(); i++ {
			thisFieldName := d.fieldNameByTag(fieldValueType.Field(i))
			thisFieldValue := fieldValue.Field(i)
			err := d.decodeStructFromNodePacket(thisFieldValue.Type(), thisFieldName, thisFieldValue, dataVal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	obtainedValue, flag := d.takeValueByKey(fieldName, fieldType, dataVal)
	if !flag {
		return fmt.Errorf("not fond fieldName:%#x", fieldName)
	}

	if !fieldValue.IsValid() {
		// This should never happen
		panic("structField is not valid")
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("fieldValue can not set:%v", fieldValue)
	}

	//fmt.Printf("#306 obtainedValue=%v, type=%v\n", obtainedValue, obtainedValue.Type())
	fieldValue.Set(obtainedValue)

	key := keyOf(fieldName)
	if !d.allowCustomizedKey(key) && key != startingToken {
		return fmt.Errorf("not allow key:%#x", key)
	}

	return nil
}

func (d *Decoder) getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	//case kind >= reflect.Int && kind <= reflect.Int64:
	//	return reflect.Int
	//case kind >= reflect.Uint && kind <= reflect.Uint64:
	//	return reflect.Uint
	//case kind >= reflect.Float32 && kind <= reflect.Float64:
	//	return reflect.Float32
	default:
		return kind
	}
}

func (d *Decoder) takeValueByKey(name string, fieldType reflect.Type, node *y3.NodePacket) (reflect.Value, bool) {
	key := keyOf(name)
	flag, isNode, packet := d.matchingKey(key, node)
	//fmt.Printf("#404 flag=%v, isNode=%v, packet_type=%v\n", flag, isNode, reflect.ValueOf(packet).Type())
	if flag == false {
		return reflect.Indirect(reflect.ValueOf(packet)), false
	}
	if isNode {
		nodePacket := packet.(y3.NodePacket)
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
		return reflect.Indirect(reflect.ValueOf(packet)), true
	}

	primitivePacket := packet.(y3.PrimitivePacket)
	return d.takePrimitiveValue(fieldType, primitivePacket)

	//switch fieldType.Kind() {
	//case reflect.String:
	//	val, err := primitivePacket.ToUTF8String()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Int32:
	//	val, err := primitivePacket.ToInt32()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Int64:
	//	val, err := primitivePacket.ToInt64()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Uint32:
	//	val, err := primitivePacket.ToUInt32()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Uint64:
	//	val, err := primitivePacket.ToUInt64()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Float32:
	//	val, err := primitivePacket.ToFloat32()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	//case reflect.Float64:
	//	val, err := primitivePacket.ToFloat64()
	//	if err != nil {
	//		panic(err)
	//	}
	//	return reflect.Indirect(reflect.ValueOf(val)), true
	////case reflect.Array:
	////	fmt.Printf("#404 Array %v\n", )
	//default:
	//	panic(errors.New("::takeValueByKey error: no matching type"))
	//}
}

func (d *Decoder) paddingToArray(fieldType reflect.Type, nodePacket y3.NodePacket) reflect.Value {
	sliceType := reflect.SliceOf(fieldType.Elem())
	sliceValue := reflect.New(sliceType).Elem()

	items := make([]reflect.Value, 0)
	for _, primitivePacket := range nodePacket.PrimitivePackets {
		itemValue, _ := d.takePrimitiveValue(fieldType.Elem(), primitivePacket)
		//fmt.Printf("#404 itemValue=%v, type=%v\n", itemValue, itemValue.Type())
		items = append(items, itemValue)
	}
	slice := reflect.Append(sliceValue, items...)

	arrayType := reflect.ArrayOf(len(nodePacket.PrimitivePackets), fieldType.Elem())
	arrayValue := reflect.New(arrayType).Elem()
	reflect.Copy(arrayValue, slice) // TODO: 注意性能影响
	return arrayValue
}

func (d *Decoder) paddingToSlice(fieldType reflect.Type, nodePacket y3.NodePacket) reflect.Value {
	sliceType := reflect.SliceOf(fieldType.Elem())
	sliceValue := reflect.New(sliceType).Elem()

	items := make([]reflect.Value, 0)
	for _, primitivePacket := range nodePacket.PrimitivePackets {
		itemValue, _ := d.takePrimitiveValue(fieldType.Elem(), primitivePacket)
		//fmt.Printf("#404 itemValue=%v, type=%v\n", itemValue, itemValue.Type())
		items = append(items, itemValue)
	}
	slice := reflect.Append(sliceValue, items...)

	return slice
}

func (d *Decoder) takePrimitiveValue(fieldType reflect.Type, primitivePacket y3.PrimitivePacket) (reflect.Value, bool) {
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
	default:
		panic(errors.New("::takeValueByKey error: no matching type"))
	}
}

func (d *Decoder) matchingKey(key byte, node *y3.NodePacket) (flag bool, isNode bool, packet interface{}) {
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

func (d *Decoder) allowCustomizedKey(key byte) bool {
	switch key {
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return false
	}
	return true
}

func (d *Decoder) fieldNameByTag(field reflect.StructField) string {
	fieldName := field.Name

	tagValue := field.Tag.Get(d.config.TagName)
	tagValue = strings.SplitN(tagValue, ",", 2)[0]
	if tagValue != "" {
		fieldName = tagValue
	}

	return fieldName
}
