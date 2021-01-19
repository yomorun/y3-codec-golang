package y3

import (
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func TestStructEncoderWithSignals(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := NewStructEncoder(0x30, StructEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))
}

func TestStructEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := NewStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestStructSliceEncoderWithSignals(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := NewStructEncoder(0x30, StructEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestStructSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := NewStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		CreateSignal(0x02).SetString("a"),
		CreateSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}
