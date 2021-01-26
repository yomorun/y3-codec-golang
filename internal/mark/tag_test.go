package mark

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNode(t *testing.T) {
	var expected byte = 0x81
	tag := NewTag(expected)
	assert.True(t, tag.IsNode(), "should be a node")
	assert.False(t, tag.IsSlice(), "It should not be a Slice")
	assert.Equal(t, expected, tag.Raw(), fmt.Sprintf("value does not match(%v): %v", expected, tag.Raw()))
	assert.Equal(t, byte(0x01), tag.SeqID(), fmt.Sprintf("value does not match(%v): %v", 0x01, tag.SeqID()))
}
