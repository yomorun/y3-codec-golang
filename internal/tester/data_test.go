package tester

import (
	"fmt"
	"testing"

	"github.com/yomorun/y3-codec-golang"

	"github.com/stretchr/testify/assert"
)

func TestBasicTestData(t *testing.T) {
	input := BasicTestData{
		Vstring:  "foo",
		Vint32:   int32(127),
		Vint64:   int64(-1),
		Vuint32:  uint32(130),
		Vuint64:  uint64(18446744073709551615),
		Vfloat32: float32(0.25),
		Vfloat64: float64(23),
		Vbool:    true,
	}
	assert.NotEmpty(t, input, "Should not equal empty")
	assert.Equal(t, "foo", input.Vstring, fmt.Sprintf("value does not match(%v): %v", "foo", input.Vstring))
}

func TestObservableTestData(t *testing.T) {
	type ObservableTestData struct {
		A float32 `y3:"0x10"`
		B string  `y3:"0x11"`
	}

	codec := y3.NewCodec(0x20)
	obj := ObservableTestData{A: float32(456), B: "y3"}
	buf, _ := codec.Marshal(obj)
	//fmt.Printf("%#v\n", buf)
	target := []byte{0x81, 0xa, 0xa0, 0x8, 0x10, 0x2, 0x43, 0xe4, 0x11, 0x2, 0x79, 0x33}
	for i, v := range target {
		assert.Equal(t, v, buf[i], fmt.Sprintf("should be: [%#x], but is [%#x]", v, buf[i]))
	}
}
