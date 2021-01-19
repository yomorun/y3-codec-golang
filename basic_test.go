package y3

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func TestBasicEncoderWithSignals(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoder(0x10, BasicEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32(inputBuf[2+3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input, value, fmt.Sprintf("value does not match(%v): %v", input, value))
}

func TestBasicEncoderWithSignalsNoRoot(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32(inputBuf[3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input, value, fmt.Sprintf("value does not match(%v): %v", input, value))
}

func TestBasicSliceEncoderWithSignals(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoder(0x10, BasicEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32Slice(inputBuf[2+3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	expectedValue := reflect.ValueOf(input)
	resultValue := reflect.ValueOf(value)
	for i := 0; i < expectedValue.Len(); i++ {
		assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
			fmt.Sprintf("Item values are not equal %v: %v",
				expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
	}
}

func TestBasicSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32Slice(inputBuf[3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	expectedValue := reflect.ValueOf(input)
	resultValue := reflect.ValueOf(value)
	for i := 0; i < expectedValue.Len(); i++ {
		assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
			fmt.Sprintf("Item values are not equal %v: %v",
				expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
	}
}
