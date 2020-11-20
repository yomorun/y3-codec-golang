package packetstructure

import (
	"errors"
	"fmt"
	"reflect"

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
	//fmt.Printf("#21 Type=%v, Kind=%v\n", reflect.TypeOf(output), reflect.TypeOf(output).Kind())
	//tmp := reflect.ValueOf(output)
	//fmt.Printf("#21 tmp Type=%v, Kind=%v\n", tmp.Elem().Type(), tmp.Elem().Kind())
	////fmt.Printf("#21 tmp Type=%v, Kind=%v\n", tmp.Elem().Elem().Type(), tmp.Elem().Elem().Kind())

	config := &DecoderConfig{
		// #1
		Result: output,
		// #2
		//Result:
	}

	//tmp2 := reflect.ValueOf(config.Result).Elem()
	//fmt.Printf("#21 tmp2 Type=%v, Kind=%v\n", tmp2.Type(), tmp2.Kind())
	////tmp3 := tmp2.Elem()
	////fmt.Printf("#21 tmp3 Type=%v, Kind=%v\n", tmp3.Type(), tmp3.Kind())

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
	case reflect.Struct:
		err = d.decodeStruct(input, outVal)
	case reflect.Slice:
		err = d.decodeSlice(input, outVal)
	default:
		return fmt.Errorf("%#x: unsupported type: %s", name, outputKind)
	}

	return err
}

func (d *Decoder) Bare(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Interface {
		return v
	}
	return v.Elem()
}

func (d *Decoder) decodeSlice(data *y3.NodePacket, val reflect.Value) error {
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

func (d *Decoder) decodeStruct(data *y3.NodePacket, val reflect.Value) error {
	var fields []field
	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		structField := valType.Field(i)
		fields = append(fields, field{structField, val.Field(i)})
	}
	//debug:
	//fmt.Printf("#202 val val=%v, Type=%v, CanSet=%v\n", val, val.Type(), val.CanSet())
	//fmt.Printf("#202 structs len=%v, fields=%v\n", len(fields), fields)
	//for _, f := range fields {
	//	fmt.Printf("#202 field Name=%v, CanSet=%v\n", f.field.Name, f.val.CanSet())
	//}

	for _, f := range fields {
		structField, fieldValue := f.field, f.val
		err := d.decodeStructFromNodePacket(structField.Type, fieldNameByTag(d.config.TagName, structField), fieldValue, data)
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
			currentFieldName := fieldNameByTag(d.config.TagName, fieldValueType.Field(i))
			currentFieldValue := fieldValue.Field(i)
			err := d.decodeStructFromNodePacket(currentFieldValue.Type(), currentFieldName, currentFieldValue, dataVal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	//debug:
	//if fieldType.Kind() == reflect.Slice {
	//	fmt.Printf("#606 fieldValue=%v\n", fieldValue)
	//}

	obtainedValue, flag := d.takeValueByKey(fieldName, fieldType, dataVal)
	if !flag {
		return fmt.Errorf("not fond fieldName:%#x", fieldName)
	}

	if !fieldValue.IsValid() {
		// This should never happen
		panic("structField is not valid")
	}

	//fmt.Printf("#40 fieldValue=%v, fieldType=%v, fieldName=%v\n", fieldValue, fieldType, fieldName)
	if !fieldValue.CanSet() {
		return fmt.Errorf("fieldValue can not set:%v", fieldValue)
	}

	fieldValue.Set(obtainedValue)

	key := keyOf(fieldName)
	if !d.allowCustomizedKey(key) && key != startingToken {
		return fmt.Errorf("not allow key:%#x", key)
	}

	return nil
}

func (d *Decoder) getKind(val reflect.Value) reflect.Kind {
	return val.Kind()
}

func (d *Decoder) takeValueByKey(fieldName string, fieldType reflect.Type, node *y3.NodePacket) (reflect.Value, bool) {
	key := keyOf(fieldName)
	flag, isNode, packet := d.matchingKey(key, node)
	if flag == false {
		return reflect.Indirect(reflect.ValueOf(packet)), false
	}
	if isNode {
		nodePacket := packet.(y3.NodePacket)
		return d.takeNodeValue(fieldType, nodePacket)
	}

	primitivePacket := packet.(y3.PrimitivePacket)
	return d.takePrimitiveValue(fieldType, primitivePacket)
}

func (d *Decoder) takeNodeValue(fieldType reflect.Type, nodePacket y3.NodePacket) (reflect.Value, bool) {
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

func (d *Decoder) paddingToStruct(fieldType reflect.Type, nodePacket y3.NodePacket) reflect.Value {
	obj := reflect.New(fieldType)
	obj = reflect.Indirect(obj)

	for _, v := range nodePacket.PrimitivePackets {
		for i := 0; i < obj.NumField(); i++ {
			structField := obj.Type().Field(i)
			valueField := obj.Field(i)
			fieldName := fieldNameByTag(d.config.TagName, structField)
			if v.SeqID() == keyOf(fieldName) {
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
			if v.SeqID() == keyOf(fieldName) {
				vv, _ := d.takeNodeValue(structField.Type, v)
				valueField.Set(vv)
			}
		}
	}

	return reflect.Indirect(obj)
}

func (d *Decoder) paddingToArray(fieldType reflect.Type, nodePacket y3.NodePacket) reflect.Value {
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

func (d *Decoder) paddingToSlice(fieldType reflect.Type, nodePacket y3.NodePacket) reflect.Value {
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