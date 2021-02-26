package common

import (
	"github.com/yomorun/y3-codec-golang/internal/utils"
	"github.com/yomorun/y3-codec-golang/pkg/encoding"
)

// DecodeLength decode to length
func DecodeLength(buf []byte) (length int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	return
}

// IsRootTag judge if it is the root node
func IsRootTag(b byte) bool {
	return b == (utils.MSB | utils.RootToken)
}
