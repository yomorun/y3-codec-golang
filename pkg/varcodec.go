package encoding

import (
	"errors"
)

// ErrBufferInsufficient describes error when encode/decode malformed VarInt
var ErrBufferInsufficient = errors.New("buffer insufficient")

// VarIntCodec for encode/decode VarInt
type VarCodec struct {
	// next ptr in buf
	Ptr int
	// Encoder: bytes are to be written
	// Decoder: bytes have been consumed
	Size int
}
