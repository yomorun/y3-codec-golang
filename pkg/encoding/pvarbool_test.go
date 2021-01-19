package encoding

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPVarBool(t *testing.T) {
	testPVarBool(t, true, []byte{0x01})
	testPVarBool(t, false, []byte{0x00})
}

func testPVarBool(t *testing.T, value bool, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%v): %X", value, value, bytes)
	var size = SizeOfPVarUInt32(uint32(1))
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodePVarBool(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val bool
	codec = VarCodec{}
	assert.Nil(t, codec.DecodePVarBool(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}
