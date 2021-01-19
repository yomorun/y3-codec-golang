package y3

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestToObject(t *testing.T) {
	input := thermometer{Temperature: float32(30), Humidity: float32(40)}
	inputBuf, _ := NewCodec(0x12).Marshal(input)

	var mold thermometer
	err := ToObject(inputBuf[2:], &mold)

	assert.NoError(t, err, fmt.Sprintf("converter error:%v", err))
	assert.Equal(t, input.Temperature, mold.Temperature, fmt.Sprintf("tester %v: %X", input.Temperature, mold.Temperature))
	assert.Equal(t, input.Humidity, mold.Humidity, fmt.Sprintf("tester %v: %X", input.Humidity, mold.Humidity))
	testPrintf("mold=%v, Temperature=%v, err=%v\n", mold, mold.Temperature, err)
}

func TestToObjectSlice(t *testing.T) {
	input := []thermometer{{Temperature: float32(30), Humidity: float32(40)}}
	inputBuf, _ := NewCodec(0x12).Marshal(input)

	var mold []thermometer
	err := ToObject(inputBuf[2:], &mold)

	assert.NoError(t, err, fmt.Sprintf("converter error:%v", err))
	assert.Equal(t, input[0].Temperature, mold[0].Temperature, fmt.Sprintf("tester %v: %X", input[0].Temperature, mold[0].Temperature))
	assert.Equal(t, input[0].Humidity, mold[0].Humidity, fmt.Sprintf("tester %v: %X", input[0].Humidity, mold[0].Humidity))
	testPrintf("mold[0].Temperature=%v, mold[0].Humidity=%v, err=%v\n", mold[0].Temperature, mold[0].Humidity, err)
}

func TestToInt32(t *testing.T) {
	expected := int32(127)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToInt32(v)
	})
}

func TestToInt32Slice(t *testing.T) {
	expected := []int32{-1, 127}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToInt32Slice(v)
	})
}

func TestToUInt32(t *testing.T) {
	expected := uint32(128)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUInt32(v)
	})
}

func TestToUInt32Slice(t *testing.T) {
	expected := []uint32{128, 130}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUInt32Slice(v)
	})
}

func TestToInt64(t *testing.T) {
	expected := int64(123)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToInt64(v)
	})
}

func TestToInt64Slice(t *testing.T) {
	expected := []int64{1, -1}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToInt64Slice(v)
	})
}

func TestToUInt64(t *testing.T) {
	expected := uint64(567)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUInt64(v)
	})
}

func TestToUInt64Slice(t *testing.T) {
	expected := []uint64{0, 256}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUInt64Slice(v)
	})
}

func TestToFloat32(t *testing.T) {
	expected := float32(456)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToFloat32(v)
	})
}

func TestFloat32Slice(t *testing.T) {
	expected := []float32{0.25, 0.375}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToFloat32Slice(v)
	})
}

func TestToFloat64(t *testing.T) {
	expected := float64(23)
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToFloat64(v)
	})
}

func TestFloat64Slice(t *testing.T) {
	expected := []float64{1, 2}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToFloat64Slice(v)
	})
}

func TestToBool(t *testing.T) {
	expected := true
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToBool(v)
	})
}

func TestBoolSlice(t *testing.T) {
	expected := []bool{true, false}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToBoolSlice(v)
	})
}

func TestToUTF8String(t *testing.T) {
	expected := "yomo"
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasic(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUTF8String(v)
	})
}

func TestToUTF8StringSlice(t *testing.T) {
	expected := []string{"iot", "edge"}
	inputBuf, _ := NewCodec(0x10).Marshal(expected)
	testBasicSlice(t, expected, inputBuf[2:], func(v []byte) (interface{}, error) {
		return ToUTF8StringSlice(v)
	})
}

func testBasic(t *testing.T, expected interface{}, buf []byte, converter func(v []byte) (interface{}, error)) {
	result, err := converter(buf)
	testPrintf("value=%v err=%v\n", result, err)

	assert.NoError(t, err, fmt.Sprintf("converter error:%v", err))
	assert.Equal(t, expected, result, fmt.Sprintf("tester %v: %X", expected, utils.FormatBytes(buf)))
}

func testBasicSlice(t *testing.T, expected interface{}, buf []byte, converter func(v []byte) (interface{}, error)) {
	result, err := converter(buf)
	testPrintf("value=%v err=%v\n", result, err)

	expectedValue := reflect.ValueOf(expected)
	resultValue := reflect.ValueOf(result)

	assert.Equal(t, expected, result, fmt.Sprintf("tester %v: %X", expected, utils.FormatBytes(buf)))
	assert.Equal(t, expectedValue.Kind(), resultValue.Kind(), fmt.Sprintf("Kind mismatch %v: %v", expectedValue.Kind(), resultValue.Kind()))
	assert.Equal(t, expectedValue.Len(), resultValue.Len(), fmt.Sprintf("Len is not equal %v: %v", expectedValue.Len(), resultValue.Len()))

	for i := 0; i < expectedValue.Len(); i++ {
		assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
			fmt.Sprintf("Item values are not equal %v: %v", expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
	}
}
