package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringSlice(t *testing.T) {
	value := reflect.ValueOf([]string{"a", "b"})
	out, ok := ToStringSlice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, "a", out[0], fmt.Sprintf("value does not match(%v): %v", "a", out[0]))
	assert.Equal(t, "b", out[1], fmt.Sprintf("value does not match(%v): %v", "a", out[1]))
}

func TestToInt64Slice(t *testing.T) {
	value := reflect.ValueOf([]int64{1, 2})
	out, ok := ToInt64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, int64(1), out[0], fmt.Sprintf("value does not match(%v): %v", int64(1), out[0]))
	assert.Equal(t, int64(2), out[1], fmt.Sprintf("value does not match(%v): %v", int64(2), out[1]))
}

func TestToUInt64Slice(t *testing.T) {
	value := reflect.ValueOf([]uint64{1, 2})
	out, ok := ToUInt64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, uint64(1), out[0], fmt.Sprintf("value does not match(%v): %v", uint64(1), out[0]))
	assert.Equal(t, uint64(2), out[1], fmt.Sprintf("value does not match(%v): %v", uint64(2), out[1]))
}

func TestToUFloat64Slice(t *testing.T) {
	value := reflect.ValueOf([]float64{1, 2})
	out, ok := ToUFloat64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, float64(1), out[0], fmt.Sprintf("value does not match(%v): %v", float64(1), out[0]))
	assert.Equal(t, float64(2), out[1], fmt.Sprintf("value does not match(%v): %v", float64(2), out[1]))
}

func TestToBoolSlice(t *testing.T) {
	value := reflect.ValueOf([]bool{true, false})
	out, ok := ToBoolSlice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, true, out[0], fmt.Sprintf("value does not match(%v): %v", true, out[0]))
	assert.Equal(t, false, out[1], fmt.Sprintf("value does not match(%v): %v", false, out[1]))
}
