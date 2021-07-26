package y3

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInt32(t *testing.T) {
	testMarshalBasic(t, int32(456), func(v []byte) (interface{}, error) {
		return ToInt32(v)
	})
}

func TestMarshalUInt32(t *testing.T) {
	testMarshalBasic(t, uint32(456), func(v []byte) (interface{}, error) {
		return ToUInt32(v)
	})
}

func TestMarshalInt64(t *testing.T) {
	testMarshalBasic(t, int64(456), func(v []byte) (interface{}, error) {
		return ToInt64(v)
	})
}

func TestMarshalUInt64(t *testing.T) {
	testMarshalBasic(t, uint64(456), func(v []byte) (interface{}, error) {
		return ToUInt64(v)
	})
}

func TestMarshalFloat32(t *testing.T) {
	testMarshalBasic(t, float32(456), func(v []byte) (interface{}, error) {
		return ToFloat32(v)
	})
}

func TestMarshalFloat64(t *testing.T) {
	testMarshalBasic(t, float64(23), func(v []byte) (interface{}, error) {
		return ToFloat64(v)
	})
}

func TestMarshalBool(t *testing.T) {
	testMarshalBasic(t, true, func(v []byte) (interface{}, error) {
		return ToBool(v)
	})
}

func TestMarshalUTF8String(t *testing.T) {
	testMarshalBasic(t, "yomo", func(v []byte) (interface{}, error) {
		return ToUTF8String(v)
	})
}

func TestMarshalEmptyUTF8String(t *testing.T) {
	testMarshalBasic(t, "", func(v []byte) (interface{}, error) {
		return ToUTF8String(v)
	})
}

func TestMarshalInt32Slice(t *testing.T) {
	testMarshalBasicSlice(t, []int32{123, 456}, func(v []byte) (interface{}, error) {
		return ToInt32Slice(v)
	})
}

func TestMarshalUInt32Slice(t *testing.T) {
	testMarshalBasicSlice(t, []uint32{11, 22}, func(v []byte) (interface{}, error) {
		return ToUInt32Slice(v)
	})
}

func TestMarshalInt64Slice(t *testing.T) {
	testMarshalBasicSlice(t, []int64{-4097, 255}, func(v []byte) (interface{}, error) {
		return ToInt64Slice(v)
	})
}

func TestMarshalUInt64Slice(t *testing.T) {
	testMarshalBasicSlice(t, []uint64{0, 1}, func(v []byte) (interface{}, error) {
		return ToUInt64Slice(v)
	})
}

func TestMarshalFloat32Slice(t *testing.T) {
	testMarshalBasicSlice(t, []float32{0.25, 0.375}, func(v []byte) (interface{}, error) {
		return ToFloat32Slice(v)
	})
}

func TestMarshalFloat64Slice(t *testing.T) {
	testMarshalBasicSlice(t, []float64{0.12, 0.45}, func(v []byte) (interface{}, error) {
		return ToFloat64Slice(v)
	})
}

func TestMarshalBoolSlice(t *testing.T) {
	testMarshalBasicSlice(t, []bool{true, false}, func(v []byte) (interface{}, error) {
		return ToBoolSlice(v)
	})
}

func TestMarshalUTF8StringSlice(t *testing.T) {
	testMarshalBasicSlice(t, []string{"a", "b"}, func(v []byte) (interface{}, error) {
		return ToUTF8StringSlice(v)
	})
}

func testMarshalBasicSlice(t *testing.T, expected interface{}, converter func(v []byte) (interface{}, error)) {
	flag := false

	input := expected
	codec := NewCodec(0x10)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x10, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		result, err := converter(v)
		testPrintf("value=%v\n", result)

		expectedValue := reflect.ValueOf(expected)
		resultValue := reflect.ValueOf(result)

		assert.Equal(t, expectedValue.Kind(), resultValue.Kind(), fmt.Sprintf("Kind mismatch %v: %v", expectedValue.Kind(), resultValue.Kind()))
		assert.Equal(t, expectedValue.Len(), resultValue.Len(), fmt.Sprintf("Len is not equal %v: %v", expectedValue.Len(), resultValue.Len()))

		for i := 0; i < expectedValue.Len(); i++ {
			assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
				fmt.Sprintf("Item values are not equal %v: %v",
					expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
		}

		return result, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

func testMarshalBasic(t *testing.T, expected interface{}, converter func(v []byte) (interface{}, error)) {
	flag := false

	input := expected
	codec := NewCodec(0x10)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%# x\n", inputBuf)

	testDecoder(0x10, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		value, err := converter(v)
		testPrintf("value=%v\n", value)
		if err != nil {
			assert.Nil(t, err, ">>>>err is not nil>>>>")
		}
		// assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
		assert.Equal(t, expected, value, fmt.Sprintf("value does not match(%v): %v", expected, value))
		return value, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

func TestMarshalObject(t *testing.T) {
	flag := false
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	codec := NewCodec(0x30)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x12, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		testPrintf("v=%#x\n", v)
		var mold thermometer
		err := ToObject(v, &mold)
		assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
		testPrintf("mold=%v\n", mold)
		assert.Equal(t, float32(30), mold.Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold.Temperature))
		assert.Equal(t, float32(40), mold.Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold.Humidity))
		return mold, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

func TestMarshalObjectSlice(t *testing.T) {
	flag := false
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	codec := NewCodec(0x30)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x12, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		testPrintf("v=%#x\n", v)
		var mold []thermometer
		err := ToObject(v, &mold)
		assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
		testPrintf("mold=%v\n", mold)
		assert.Equal(t, float32(30), mold[0].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold[0].Temperature))
		assert.Equal(t, float32(40), mold[0].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold[0].Humidity))
		assert.Equal(t, float32(50), mold[1].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(50), mold[1].Temperature))
		assert.Equal(t, float32(60), mold[1].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(60), mold[1].Humidity))
		return mold, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

type exampleData struct {
	Name  string      `y3:"0x10"`
	Noise float32     `y3:"0x11"`
	Therm thermometer `y3:"0x12"`
}

type thermometer struct {
	Temperature float32 `y3:"0x13"`
	Humidity    float32 `y3:"0x14"`
}

type exampleSlice struct {
	Therms []thermometer `y3:"0x12"`
}
