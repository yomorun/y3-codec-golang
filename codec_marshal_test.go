package y3

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	testMarshalPrimitiveType(t, []byte{0x79, 0x2d, 0x6e, 0x65, 0x77}, "y-new")
	testMarshalPrimitiveType(t, []byte{0x02}, int32(2))
	testMarshalPrimitiveType(t, []byte{0x02}, uint32(2))
	testMarshalPrimitiveType(t, []byte{0x0}, int64(0))
	testMarshalPrimitiveType(t, []byte{0x02}, uint64(2))
	testMarshalPrimitiveType(t, []byte{0x3e, 0xc0}, float32(0.375))
	testMarshalPrimitiveType(t, []byte{0x40, 0x37}, float64(23))

	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x1, 0x61, 0x0, 0x1, 0x62, 0x0, 0x1, 0x63}, "a", "b", "c")
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x7f}, int32(1), int32(2), int32(-1))
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x3}, uint32(1), uint32(2), uint32(3))
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x7f}, int64(1), int64(2), int64(-1))
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x3}, uint64(1), uint64(2), uint64(3))
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x2, 0x3e, 0x80, 0x0, 0x2, 0x3e, 0xc0}, float32(0.25), float32(0.375))
	testMarshalPrimitiveTypeArray(t, []byte{0x0, 0x2, 0x3f, 0xf0, 0x0, 0x1, 0xc0}, float64(1), float64(-2))
}

func testMarshalPrimitiveType(t *testing.T, expected []byte, T interface{}) {
	var msg = fmt.Sprintf("testing %v, (%v)", expected, T)
	codec := NewCodec("")
	buf, _ := codec.Marshal(T)
	assert.True(t, bytes.Equal(expected, buf), msg)
}

func testMarshalPrimitiveTypeArray(t *testing.T, expected []byte, T ...interface{}) {
	var msg = fmt.Sprintf("testing %v, (%v)", expected, T)
	codec := NewCodec("")
	buf, _ := codec.Marshal(T)
	assert.True(t, bytes.Equal(expected, buf), msg)
}
