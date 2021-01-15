package y3

import (
	"fmt"
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

func TestMarshalUTF8String(t *testing.T) {
	testMarshalBasic(t, "yomo", func(v []byte) (interface{}, error) {
		return ToUTF8String(v)
	})
}

func TestMarshalFloat32Slice(t *testing.T) {
	flag := false
	expected := []float32{0.25, 0.375}

	var input interface{} = expected
	codec := NewCodec(0x10)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x10, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		value, err := ToFloat32Slice(v)
		testPrintf("value=%v\n", value)
		assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
		//assert.Equal(t, expected, value, fmt.Sprintf("value does not match(%v): %v", expected, value))
		return value, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

func testMarshalBasicSlice(t *testing.T, expected interface{}, converter func(v []byte) (interface{}, error)) {

}

func testMarshalBasic(t *testing.T, expected interface{}, converter func(v []byte) (interface{}, error)) {
	flag := false

	input := expected
	codec := NewCodec(0x10)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x10, inputBuf, func(v []byte) (interface{}, error) {
		flag = true
		value, err := converter(v)
		testPrintf("value=%v\n", value)
		assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
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

		return mold, err
	})

	if !flag {
		t.Errorf("Observable does not listen to values")
	}
}

func TestMarshal(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	codec := NewCodec(0x30)
	inputBuf, _ := codec.Marshal(input)
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	testDecoder(0x11, inputBuf, func(v []byte) (interface{}, error) {
		value, err := ToFloat32(v)
		expected := float32(456)
		assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
		assert.Equal(t, expected, value, fmt.Sprintf("value does not match(%v): %v", expected, value))
		return value, err
	})
}

type exampleData struct {
	Name  string      `yomo:"0x10"`
	Noise float32     `yomo:"0x11"`
	Therm thermometer `yomo:"0x12"`
}

type thermometer struct {
	Temperature float32 `yomo:"0x13"`
	Humidity    float32 `yomo:"0x14"`
}
