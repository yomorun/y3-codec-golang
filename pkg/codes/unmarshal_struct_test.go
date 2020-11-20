package codes

import (
	"testing"
)

type Thermometer struct {
	Id          string  `yomo:"0x10"`
	Temperature float32 `yomo:"0x11"`
	Humidity    float32 `yomo:"0x12"`
	Stored      bool    `yomo:"0x13"`
}

func TestUnmarshalStructThermometerSlice(t *testing.T) {
	input := []Thermometer{
		{
			Id:          "the0",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
		{
			Id:          "the1",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
	}

	codec1 := NewCodec("0x20")
	inputBuf, _ := codec1.Marshal(input)
	//fmt.Printf("#60 buf=%s\n", FormatBytes(inputBuf))

	codec2 := NewCodec("0x20")
	var mold []Thermometer
	runUnmarshalStruct(t, codec2, inputBuf, &mold)
	//fmt.Printf("#60 mold=%v\n", mold)

	if len(mold) != len(input) {
		t.Errorf("len value should be: %v", len(input))
	}

	for i, v := range mold {
		testThermometerStruct(t, v, input[i])
	}
}

func TestUnmarshalStructThermometer(t *testing.T) {
	input := Thermometer{
		Id:          "the0",
		Temperature: float32(64.88),
		Humidity:    float32(93.02),
		Stored:      true,
	}

	codec1 := NewCodec("0x20")
	inputBuf, _ := codec1.Marshal(input)
	//fmt.Printf("#60 buf=%s\n", FormatBytes(inputBuf))

	codec2 := NewCodec("0x20")
	var mold = Thermometer{}
	runUnmarshalStruct(t, codec2, inputBuf, &mold)
	//fmt.Printf("#60 mold=%v\n", mold)
	testThermometerStruct(t, mold, input)
}

func runUnmarshalStruct(t *testing.T, codec YomoCodec, inputBuf []byte, mold interface{}) {
	err := codec.UnmarshalStruct(inputBuf, mold)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}
}

func testThermometerStruct(t *testing.T, result Thermometer, expected Thermometer) {
	if result.Id != expected.Id {
		t.Errorf("Id value should be: %v", expected.Id)
	}

	if result.Temperature != expected.Temperature {
		t.Errorf("Temperature value should be: %v", expected.Temperature)
	}

	if result.Humidity != expected.Humidity {
		t.Errorf("Humidity value should be: %v", expected.Humidity)
	}

	if result.Stored != expected.Stored {
		t.Errorf("Stored value should be: %v", expected.Stored)
	}
}
