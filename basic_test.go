package y3

import (
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func TestBasicEncoderWithSignals(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoder(0x10, BasicEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))
}

func TestBasicEncoderWithSignalsNoRoot(t *testing.T) {
	input := int32(456)

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestBasicSliceEncoderWithSignals(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoder(0x10, BasicEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestBasicSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := []int32{123, 456}

	encoder := NewBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))
}
