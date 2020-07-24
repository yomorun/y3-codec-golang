package encoding

import (
	"errors"
)

var BufferInsufficient = errors.New("buffer insufficient")

type VarIntCodec struct {
	// next ptr in buf
	Ptr  int
	// Encoder: bytes are to be written
	// Decoder: bytes have been consumed
	Size int
}
