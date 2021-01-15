package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func ToSliceArray(arg interface{}) (arr []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		return value.Interface()
	})
}

func ToSliceArrayWith(arg interface{}, handle func(value reflect.Value) interface{}) (arr []interface{}, ok bool) {
	argValue := reflect.ValueOf(arg)
	if argValue.Type().Kind() == reflect.Slice {
		length := argValue.Len()
		if length == 0 {
			return
		}
		ok = true
		arr = make([]interface{}, length)
		for i := 0; i < length; i++ {
			arr[i] = handle(argValue.Index(i))
		}
	}
	return
}

func ToStringSliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		return fmt.Sprintf("%v", value)
	})
}

func ToInt64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		i64, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
		return i64
	})
}

func ToUInt64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		ui64, _ := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)
		return ui64
	})
}

func ToUFloat64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		f64, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		return f64
	})
}

func ToBoolSliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		bl, _ := strconv.ParseBool(fmt.Sprintf("%v", value))
		return bl
	})
}
