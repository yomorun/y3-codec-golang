package y3

import (
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func TestBasicEncoderWithSignals(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoderWithRoot(0x10, rootToken)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestBasicEncoderWithSignalsNoRoot(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestBasicSliceEncoderWithSignals(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoderWithRoot(0x10, rootToken)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestBasicSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}
