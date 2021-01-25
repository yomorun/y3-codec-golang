package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyOf(t *testing.T) {
	value := KeyOf("0x01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))

	value = KeyOf("0X01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))

	value = KeyOf("01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))
}

func TestIsEmptyKey(t *testing.T) {
	assert.True(t, IsEmptyKey(0x00), "0x00 is a empty key")
}

func TestForbiddenCustomizedKey(t *testing.T) {
	assert.True(t, ForbiddenCustomizedKey(0x01), "0x01 is disabled")
	assert.False(t, ForbiddenCustomizedKey(0x10), "0x10 is allowed")
}

func TestAllowableSignalKey(t *testing.T) {
	assert.True(t, AllowableSignalKey(0x02), "0x01 is allowed")
	assert.False(t, AllowableSignalKey(0x01), "0x10 is disabled")
}
