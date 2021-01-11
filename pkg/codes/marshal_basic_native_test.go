package codes

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/yomorun/y3-codec-golang/pkg/packetutils"

	"github.com/stretchr/testify/assert"
)

func TestMarshalBasicNative(t *testing.T) {
	testMarshalBasicNative(t, []byte{0x79, 0x2d, 0x6e, 0x65, 0x77}, "y-new")
	testMarshalBasicNative(t, []byte{0x02}, int32(2))
	testMarshalBasicNative(t, []byte{0x02}, uint32(2))
	testMarshalBasicNative(t, []byte{0x0}, int64(0))
	testMarshalBasicNative(t, []byte{0x02}, uint64(2))
	testMarshalBasicNative(t, []byte{0x3e, 0xc0}, float32(0.375))
	testMarshalBasicNative(t, []byte{0x40, 0x37}, float64(23))

	testMarshalBasicNativeSlice(t, []byte{0x0, 0x1, 0x61, 0x0, 0x1, 0x62, 0x0, 0x1, 0x63}, []string{"a", "b", "c"})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x7f}, []int32{int32(1), int32(2), int32(-1)})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x3}, []uint32{uint32(1), uint32(2), uint32(3)})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x7f}, []int64{int64(1), int64(2), int64(-1)})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x1, 0x1, 0x0, 0x1, 0x2, 0x0, 0x1, 0x3}, []uint64{uint64(1), uint64(2), uint64(3)})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x2, 0x3e, 0x80, 0x0, 0x2, 0x3e, 0xc0}, []float32{float32(0.25), float32(0.375)})
	testMarshalBasicNativeSlice(t, []byte{0x0, 0x2, 0x3f, 0xf0, 0x0, 0x1, 0xc0}, []float64{float64(1), float64(-2)})
}

func testMarshalBasicNative(t *testing.T, expected []byte, T interface{}) {
	var msg = fmt.Sprintf("testing %v, (%v)", expected, T)
	proto := NewProtoCodec(packetutils.KeyOf(""))
	buf, _ := proto.MarshalNative(T)
	//fmt.Printf("#000 testMarshalBasicNative Kind=%v, buf=%v\n", reflect.ValueOf(T).Kind(), packetutils.FormatBytes(buf))
	assert.True(t, bytes.Equal(expected, buf), msg)
}

func testMarshalBasicNativeSlice(t *testing.T, expected []byte, T interface{}) {
	var msg = fmt.Sprintf("testing %v, (%v)", expected, T)
	proto := NewProtoCodec(packetutils.KeyOf(""))
	buf, _ := proto.MarshalNative(T)
	//fmt.Printf("#000 testMarshalBasicNativeSlice Kind=%v, Elem=%v, buf=%v\n", reflect.ValueOf(T).Kind(), reflect.ValueOf(T).Index(0).Kind(), packetutils.FormatBytes(buf))
	assert.True(t, bytes.Equal(expected, buf), msg)
}
