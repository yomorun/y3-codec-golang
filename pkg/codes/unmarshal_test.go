package codes

import (
	"fmt"
	"testing"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	"github.com/stretchr/testify/assert"
)

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
	testUnmarshalUint32(t, data, "0x12", uint32(1))
	testUnmarshalInt64(t, data, "0x13", int64(-1))
	testUnmarshalUint64(t, data, "0x15", uint64(1))
	testUnmarshalFloat32(t, data, "0x17", float32(0.25))
	testUnmarshalFloat64(t, data, "0x19", float64(1))

	testUnmarshalStringArray(t, data, "0x25", []string{"a", "b"})
	testUnmarshalInt32Array(t, data, "0x10", []int32{1, 2})
	testUnmarshalUint32Array(t, data, "0x11", []uint32{1, 2})
	testUnmarshalInt64Array(t, data, "0x14", []int64{1, 2})
	testUnmarshalUint64Array(t, data, "0x16", []uint64{1, 2})
	testUnmarshalFloat32Array(t, data, "0x18", []float32{0.25})
	testUnmarshalFloat64Array(t, data, "0x1a", []float64{1})
}

func testUnmarshalString(t *testing.T, data []byte, observe string, expected string) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	codec := NewCodec(observe)
	var mold interface{} = ""
	_ = codec.Unmarshal(data, &mold)
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
