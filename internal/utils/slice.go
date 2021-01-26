package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

// ToStringSlice converting interface to interface slice for string
func ToStringSlice(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceWith(arg, func(value reflect.Value) interface{} {
		return fmt.Sprintf("%v", value)
	})
}

// ToInt64Slice converting interface to interface slice for int64
func ToInt64Slice(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceWith(arg, func(value reflect.Value) interface{} {
		i64, _ := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
		return i64
	})
}

// ToUInt64Slice converting interface to interface slice for uint64
func ToUInt64Slice(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceWith(arg, func(value reflect.Value) interface{} {
		ui64, _ := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 64)
		return ui64
	})
}

// ToUFloat64Slice converting interface to interface slice for float64
func ToUFloat64Slice(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceWith(arg, func(value reflect.Value) interface{} {
		f64, _ := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		return f64
	})
}

// ToUFloat64Slice converting interface to interface slice for bool
func ToBoolSlice(arg interface{}) (out []interface{}, ok bool) {
	return ToSliceWith(arg, func(value reflect.Value) interface{} {
		bl, _ := strconv.ParseBool(fmt.Sprintf("%v", value))
		return bl
	})
}

// ToSliceWith converting interface to interface slice, and using a custom handle
func ToSliceWith(arg interface{}, handle func(value reflect.Value) interface{}) (arr []interface{}, ok bool) {
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
