package y3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func TestStructEncoderWithSignals(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := newStructEncoder(0x30, structEncoderOptionRoot(rootToken),
		structEncoderOptionConfig(&structEncoderConfig{
			ZeroFields: true,
			TagName:    "y3",
		}))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleData
	err := ToObject(inputBuf[2+3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input.Name, mold.Name, fmt.Sprintf("value does not match(%v): %v", input.Name, mold.Name))
	assert.Equal(t, input.Noise, mold.Noise, fmt.Sprintf("value does not match(%v): %v", input.Noise, mold.Noise))
	assert.Equal(t, input.Therm.Temperature, mold.Therm.Temperature, fmt.Sprintf("value does not match(%v): %v", input.Therm.Temperature, mold.Therm.Temperature))
	assert.Equal(t, input.Therm.Humidity, mold.Therm.Humidity, fmt.Sprintf("value does not match(%v): %v", input.Therm.Humidity, mold.Therm.Humidity))
}

func TestStructEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := newStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleData
	err := ToObject(inputBuf[3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input.Name, mold.Name, fmt.Sprintf("value does not match(%v): %v", input.Name, mold.Name))
	assert.Equal(t, input.Noise, mold.Noise, fmt.Sprintf("value does not match(%v): %v", input.Noise, mold.Noise))
	assert.Equal(t, input.Therm.Temperature, mold.Therm.Temperature, fmt.Sprintf("value does not match(%v): %v", input.Therm.Temperature, mold.Therm.Temperature))
	assert.Equal(t, input.Therm.Humidity, mold.Therm.Humidity, fmt.Sprintf("value does not match(%v): %v", input.Therm.Humidity, mold.Therm.Humidity))
}

func TestStructSliceEncoderWithSignals(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := newStructEncoder(0x30, structEncoderOptionRoot(rootToken))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleSlice
	err := ToObject(inputBuf[2+3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	assert.Equal(t, float32(30), mold.Therms[0].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold.Therms[0].Temperature))
	assert.Equal(t, float32(40), mold.Therms[0].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold.Therms[0].Humidity))
	assert.Equal(t, float32(50), mold.Therms[1].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(50), mold.Therms[1].Temperature))
	assert.Equal(t, float32(60), mold.Therms[1].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(60), mold.Therms[1].Humidity))
}

func TestStructSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := newStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleSlice
	err := ToObject(inputBuf[3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	assert.Equal(t, float32(30), mold.Therms[0].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold.Therms[0].Temperature))
	assert.Equal(t, float32(40), mold.Therms[0].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold.Therms[0].Humidity))
	assert.Equal(t, float32(50), mold.Therms[1].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(50), mold.Therms[1].Temperature))
	assert.Equal(t, float32(60), mold.Therms[1].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(60), mold.Therms[1].Humidity))
}
