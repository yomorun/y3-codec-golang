package codes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

func TestUnmarshalBasic(t *testing.T) {
	data := []byte{0x81, 0x80, 0x5f, 0x22, 0x1, 0x1, 0x12, 0x1, 0x1,
		0x13, 0x1, 0x7f, 0x15, 0x1, 0x1, 0x17, 0x2, 0x3e, 0x80, 0x23, 0x1,
		0x79, 0xa4, 0x34, 0xd4, 0x6, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd6, 0x6,
		0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd8, 0x4, 0x0, 0x2, 0x3e, 0x80, 0xda,
		0x4, 0x0, 0x2, 0x3f, 0xf0, 0xe5, 0x6, 0x0, 0x1, 0x61, 0x0, 0x1, 0x62,
		0xd0, 0x6, 0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0xd1, 0x6, 0x0, 0x1, 0x1, 0x0,
		0x1, 0x2, 0x19, 0x2, 0x3f, 0xf0, 0xe6, 0x10, 0x80, 0x6, 0x27, 0x1, 0x5,
		0x28, 0x1, 0x2, 0x80, 0x6, 0x27, 0x1, 0x6, 0x28, 0x1, 0x3}

	testUnmarshalBasicString(t, data, "0x23", "y")
	testUnmarshalBasicInt32(t, data, "0x22", int32(1))
	testUnmarshalBasicUint32(t, data, "0x12", uint32(1))
	testUnmarshalBasicInt64(t, data, "0x13", int64(-1))
	testUnmarshalBasicUint64(t, data, "0x15", uint64(1))
	testUnmarshalBasicFloat32(t, data, "0x17", float32(0.25))
	testUnmarshalBasicFloat64(t, data, "0x19", float64(1))

	testUnmarshalBasicStringArray(t, data, "0x25", []string{"a", "b"})
	testUnmarshalBasicInt32Array(t, data, "0x10", []int32{1, 2})
	testUnmarshalBasicUint32Array(t, data, "0x11", []uint32{1, 2})
	testUnmarshalBasicInt64Array(t, data, "0x14", []int64{1, 2})
	testUnmarshalBasicUint64Array(t, data, "0x16", []uint64{1, 2})
	testUnmarshalBasicFloat32Array(t, data, "0x18", []float32{0.25})
	testUnmarshalBasicFloat64Array(t, data, "0x1a", []float64{1})
}

func runUnmarshalBasic(t *testing.T, proto ProtoCodec, data []byte, mold *interface{}) {
	err := proto.UnmarshalBasic(data, mold)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}
}

func testUnmarshalBasicString(t *testing.T, data []byte, observe string, expected string) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = ""
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(string), msg)
}

func testUnmarshalBasicInt32(t *testing.T, data []byte, observe string, expected int32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = int32(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(int32), msg)
}

func testUnmarshalBasicUint32(t *testing.T, data []byte, observe string, expected uint32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = uint32(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(uint32), msg)
}

func testUnmarshalBasicInt64(t *testing.T, data []byte, observe string, expected int64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = int64(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(int64), msg)
}

func testUnmarshalBasicUint64(t *testing.T, data []byte, observe string, expected uint64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = uint64(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(uint64), msg)
}

func testUnmarshalBasicFloat32(t *testing.T, data []byte, observe string, expected float32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = float32(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(float32), msg)
}

func testUnmarshalBasicFloat64(t *testing.T, data []byte, observe string, expected float64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = float64(0)
	runUnmarshalBasic(t, proto, data, &mold)
	assert.Equal(t, expected, mold.(float64), msg)
}

func testUnmarshalBasicStringArray(t *testing.T, data []byte, observe string, expected []string) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]string{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(string), msg)
		}
	}
}

func testUnmarshalBasicInt32Array(t *testing.T, data []byte, observe string, expected []int32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]int32{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(int32), msg)
		}
	}
}

func testUnmarshalBasicUint32Array(t *testing.T, data []byte, observe string, expected []uint32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]uint32{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(uint32), msg)
		}
	}
}

func testUnmarshalBasicInt64Array(t *testing.T, data []byte, observe string, expected []int64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]int64{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(int64), msg)
		}
	}
}

func testUnmarshalBasicUint64Array(t *testing.T, data []byte, observe string, expected []uint64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]uint64{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(uint64), msg)
		}
	}
}

func testUnmarshalBasicFloat32Array(t *testing.T, data []byte, observe string, expected []float32) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]float32{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(float32), msg)
		}
	}
}

func testUnmarshalBasicFloat64Array(t *testing.T, data []byte, observe string, expected []float64) {
	var msg = fmt.Sprintf("testing %s,  %v, (%X)", observe, expected, data)
	proto := NewProtoCodec(observe)
	var mold interface{} = [0]float64{}
	runUnmarshalBasic(t, proto, data, &mold)

	if arr, ok := utils.ToSliceArray(mold); ok {
		assert.Equal(t, len(expected), len(arr), msg)
		for i, v := range arr {
			assert.Equal(t, expected[i], v.(float64), msg)
		}
	}
}
