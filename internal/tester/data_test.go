package tester

import (
	"fmt"
	"testing"

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
