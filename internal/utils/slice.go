package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

// ToSliceArray converting interface to interface slice
func ToSliceArray(arg interface{}) (arr []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		return value.Interface()
	})
}

// ToStringSliceArray converting interface to interface slice for string
func ToStringSliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		return fmt.Sprintf("%v", value)
	})
}

// ToInt64SliceArray converting interface to interface slice for int64
func ToInt64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		i64, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
		return i64
	})
}

// ToUInt64SliceArray converting interface to interface slice for uint64
func ToUInt64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		ui64, _ := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)
		return ui64
	})
}

// ToUFloat64SliceArray converting interface to interface slice for float64
func ToUFloat64SliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		f64, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		return f64
	})
}

// ToUFloat64SliceArray converting interface to interface slice for bool
func ToBoolSliceArray(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceArrayWith(arg, func(value reflect.Value) interface{} {
		bl, _ := strconv.ParseBool(fmt.Sprintf("%v", value))
		return bl
	})
}

// ToSliceArrayWith converting interface to interface slice, and using a custom handle
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
