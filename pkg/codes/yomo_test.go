package codes

import (
	"fmt"
	"testing"
)

func TestReadThermometer(t *testing.T) {
	inputBuf := buildThermometerInputData()

	codec := NewCodec("0x20")
	codec.Decoder(inputBuf)
	fmt.Printf("#71 buf=%s\n", FormatBytes(inputBuf))

	// #1
	var mold = Thermometer{}
	//var mold = getThermometerInterface()
	val, err := codec.Read(&mold)
	if err != nil {
		fmt.Printf("#71 err1=%v\n", err)
		_ = val
	}
	fmt.Printf("#71 mold=%v\n", mold)

	result, err := handleThermometer(mold)
	if err != nil {
		fmt.Printf("#71 err2=%v\n", err)
	}
	fmt.Printf("#71 result=%v\n", result)

	// #2
	//var mold = getThermometerInterface()
	//val, err := codec.Read(&mold)
	//if err != nil {
	//	fmt.Printf("#71 err1=%v\n", err)
	//	_ = val
	//}
	//fmt.Printf("#71 mold=%v\n", mold)
}

func buildThermometerInputData() []byte {
	input := Thermometer{
		Id:          "the0",
		Temperature: float32(64.88),
		Humidity:    float32(93.02),
		Stored:      true,
	}

	codec := NewCodec("0x20")
	inputBuf, _ := codec.Marshal(input)
	//fmt.Printf("#70 buf=%s\n", FormatBytes(inputBuf))
	return inputBuf
}

func handleThermometer(value interface{}) (interface{}, error) {
	the := value.(Thermometer)
	the.Id = "the1"
	return the, nil
}

func getThermometerInterface() interface{} {
	return Thermometer{}
}

//
func TestReadThermometerSlice(t *testing.T) {
	inputBuf := buildThermometerSliceInputData()

	codec := NewCodec("0x20")
	codec.Decoder(inputBuf)
	fmt.Printf("#70 buf=%s\n", FormatBytes(inputBuf))

	// #1
	var mold = []Thermometer{}
	val, err := codec.Read(&mold)
	if err != nil {
		fmt.Printf("#70 err1=%v\n", err)
		_ = val
	}
	fmt.Printf("#70 mold=%v\n", mold)
	////fmt.Printf("#70 val=%v\n", val)
	//fmt.Printf("#70 mold=%v\n", mold)
	//val = val.([]Thermometer)
	//fmt.Printf("#70 val=%v\n", val)
	result, err := handleThermometerSlice(mold)
	if err != nil {
		fmt.Printf("#70 err2=%v\n", err)
	}
	fmt.Printf("#70 result=%v\n", result)

	// #2
	//mold := getThermometerSliceInterface()
	//val, err := codec.Read(mold)
	//fmt.Printf("#70 err=%v\n", err)
	//fmt.Printf("#70 val=%v\n", val)

}

func getThermometerSliceInterface() interface{} {
	return []Thermometer{}
}

func handleThermometerSlice(value interface{}) (interface{}, error) {
	the := value.([]Thermometer)
	the = append(the, Thermometer{
		Id:          "the1",
		Temperature: float32(1),
		Humidity:    float32(2),
		Stored:      false,
	})
	return the, nil
}

func buildThermometerSliceInputData() []byte {
	input := []Thermometer{
		{
			Id:          "the0",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
		{
			Id:          "the1",
			Temperature: float32(50),
			Humidity:    float32(90),
			Stored:      true,
		},
	}

	codec := NewCodec("0x20")
	inputBuf, _ := codec.Marshal(input)
	//fmt.Printf("#70 buf=%s\n", FormatBytes(inputBuf))
	return inputBuf
}
