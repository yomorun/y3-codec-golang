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

	encoder := NewStructEncoderWithRoot(0x30, input, rootToken)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))
}

func TestStructEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := NewStructEncoder(0x30, input)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestStructSliceEncoderWithSignals(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := NewStructEncoderWithRoot(0x30, input, 0x01)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}

func TestStructSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := NewStructEncoder(0x30, input)
	inputBuf, _ := encoder.EncodeWithSignals(input, func() []*PrimitivePacketEncoder {
		from := NewPrimitivePacketEncoder(0x02)
		from.SetStringValue("a")
		to := NewPrimitivePacketEncoder(0x03)
		to.SetStringValue("b")
		return []*PrimitivePacketEncoder{from, to}
	})
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

}
