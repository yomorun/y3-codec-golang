package codes

import (
	"io"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

var (
	placeholder = []byte{0, 1, 2, 3}
	logger      = utils.Logger.WithPrefix(utils.DefaultLogger, "collectingCodec")
)

// YomoCodec: codec interface for yomo
type YomoCodec interface {
	Decoder(buf []byte)
	Read(mold interface{}) (interface{}, error)
	Write(w io.Writer, T interface{}, mold interface{}) (int, error)
	Refresh(w io.Writer) (int, error)
}

func NewCodec(observe string) YomoCodec {
	//return NewCollectingCodec(observe)
	return NewStreamingCodec(packetutils.KeyOf(observe))
}
