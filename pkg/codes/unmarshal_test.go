package codes

import (
	"fmt"
	"testing"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	"github.com/stretchr/testify/assert"
)

//type Thermometer struct {
//	Id          string  `yomo:"0x10"`
//	Temperature float32 `yomo:"0x11"`
//	Humidity    float32 `yomo:"0x12"`
//	Stored      bool    `yomo:"0x13"`
//}

//func TestUnmarshalStruct(t *testing.T) {
//	input := Thermometer{
//		Id:          "the0",
//		Temperature: float32(64.88),
//		Humidity:    float32(93.02),
//		Stored:      true,
//	}
//
//	codec1 := NewCodec("0x20")
//	inputBuf, _ := codec1.Marshal(input)
//	fmt.Printf("#30 buf=%s\n", FormatBytes(inputBuf))
//
//	codec2 := NewCodec("0x20")
//	//mold := new(Thermometer)
//	//var mold Thermometer
//	mold := Thermometer{}
//	var data interface{} = mold
//	////data := &Thermometer{}
//
//	fmt.Printf("#30 mold-type=%s, mold-kind=%v\n", reflect.ValueOf(mold).Type(), reflect.ValueOf(mold).Type().Kind())
//	fmt.Printf("#30 mold-type=%s, mold-kind=%v\n", reflect.ValueOf(&mold).Type(), reflect.ValueOf(&mold).Type().Kind())
//	//fmt.Printf("#30 data-type=%s, data-kind=%v\n", reflect.ValueOf(&data).Type(), reflect.ValueOf(&data).Type().Kind())
//
//	//err := codec2.Unmarshal2(inputBuf, &mold)
//	//fmt.Printf("#31 err=%v\n", err)
//	//fmt.Printf("#31 mold=%v\n", mold)
//
//	//err := codec2.Unmarshal(inputBuf, &data)
//	//fmt.Printf("#31 err=%v\n", err)
//	//fmt.Printf("#31 data=%v\n", data.(interface{}))
//	//fmt.Printf("#31 data=%v, Temperature=%v\n", data, data.(Thermometer).Temperature)
//
//	//Unmarshal3
//	_, err := codec2.Unmarshal3(inputBuf, &data)
//	//fmt.Printf("#31 err=%v\n", err)
//	//fmt.Printf("#31 data=%v, Temperature=%v\n", result, result.(*Thermometer).Temperature)
//	//fmt.Printf("#31 result=%v\n", result)
//
//	fmt.Printf("#31 err=%v\n", err)
//	fmt.Printf("#31 data=%v\n", data)
//	//fmt.Printf("#31 data=%v, Temperature=%v\n", mold, mold.(*Thermometer).Temperature)
//}
//
////func testA(codec YomoCodec, data []byte, v interface{}) error {
////	return codec.Unmarshal(data, &v)
////}

func TestUnmarshal(t *testing.T) {
	data := []byte{0x81, 0x80, 0x5f, 0x22, 0x1, 0x1, 0x12, 0x1, 0x1,
		0x13, 0x1, 0x7f, 0x15, 0x1, 0x1, 0x17, 0x2, 0x3e, 0x80, 0x23, 0x1,
		0x79, 0xa4, 0x34, 0xd4, 0x6, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd6, 0x6,
		0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd8, 0x4, 0x0, 0x2, 0x3e, 0x80, 0xda,
		0x4, 0x0, 0x2, 0x3f, 0xf0, 0xe5, 0x6, 0x0, 0x1, 0x61, 0x0, 0x1, 0x62,
		0xd0, 0x6, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd1, 0x6, 0x0, 0x1, 0x1, 0x0,
		0x1, 0x2, 0x19, 0x2, 0x3f, 0xf0, 0xe6, 0x10, 0x80, 0x6, 0x27, 0x1, 0x5,
		0x28, 0x1, 0x2, 0x80, 0x6, 0x27, 0x1, 0x6, 0x28, 0x1, 0x3}

	testUnmarshalString(t, data, "0x23", "y")
	testUnmarshalInt32(t, data, "0x22", int32(1))
	//testUnmarshalUint32(t, data, "0x12", uint32(1))
	//testUnmarshalInt64(t, data, "0x13", int64(-1))
	//testUnmarshalUint64(t, data, "0x15", uint64(1))
	//testUnmarshalFloat32(t, data, "0x17", float32(0.25))
	//testUnmarshalFloat64(t, data, "0x19", float64(1))
	//
	//testUnmarshalStringArray(t, data, "0x25", []string{"a", "b"})
	//testUnmarshalInt32Array(t, data, "0x10", []int32{1, 2})
	//testUnmarshalUint32Array(t, data, "0x11", []uint32{1, 2})
	//testUnmarshalInt64Array(t, data, "0x14", []int64{1, 2})
	//testUnmarshalUint64Array(t, data, "0x16", []uint64{1, 2})
	//testUnmarshalFloat32Array(t, data, "0x18", []float32{0.25})
	//testUnmarshalFloat64Array(t, data, "0x1a", []float64{1})
}

//func testUnmarshalString(t *testing.T, data []byte, observe string, expected string) {
//	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
//	codec := NewCodec(observe)
//	var mold interface{} = ""
//	_ = codec.Unmarshal(data, &mold)
//	assert.Equal(t, expected, mold.(string), msg)
//}

func testUnmarshalString(t *testing.T, data []byte, observe string, expected string) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = ""
	_ = codec.Unmarshal(data, &mold)
	fmt.Printf("#44 mold=%v\n", mold)
	assert.Equal(t, expected, mold.(string), msg)
}

func testUnmarshalInt32(t *testing.T, data []byte, observe string, expected int32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = int32(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(int32), msg)
}

func testUnmarshalUint32(t *testing.T, data []byte, observe string, expected uint32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = uint32(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(uint32), msg)
}

func testUnmarshalInt64(t *testing.T, data []byte, observe string, expected int64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = int64(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(int64), msg)
}

func testUnmarshalUint64(t *testing.T, data []byte, observe string, expected uint64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = uint64(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(uint64), msg)
}

func testUnmarshalFloat32(t *testing.T, data []byte, observe string, expected float32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = float32(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(float32), msg)
}

func testUnmarshalFloat64(t *testing.T, data []byte, observe string, expected float64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = float64(0)
	_ = codec.Unmarshal(data, &mold)
	assert.Equal(t, expected, mold.(float64), msg)
}

func testUnmarshalStringArray(t *testing.T, data []byte, observe string, expected []string) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]string{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(string), msg)
		}
	}
}

func testUnmarshalInt32Array(t *testing.T, data []byte, observe string, expected []int32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]int32{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(int32), msg)
		}
	}
}

func testUnmarshalUint32Array(t *testing.T, data []byte, observe string, expected []uint32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]uint32{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(uint32), msg)
		}
	}
}

func testUnmarshalInt64Array(t *testing.T, data []byte, observe string, expected []int64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]int64{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(int64), msg)
		}
	}
}

func testUnmarshalUint64Array(t *testing.T, data []byte, observe string, expected []uint64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]uint64{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(uint64), msg)
		}
	}
}

func testUnmarshalFloat32Array(t *testing.T, data []byte, observe string, expected []float32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]float32{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(float32), msg)
		}
	}
}

func testUnmarshalFloat64Array(t *testing.T, data []byte, observe string, expected []float64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = [0]float64{}
	_ = codec.Unmarshal(data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(float64), msg)
		}
	}
}
