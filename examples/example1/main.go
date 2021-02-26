package main

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/yomorun/y3-codec-golang"
)

// SourceData define struct for this example
type SourceData struct {
	Name  string      `y3:"0x10"`
	Noise float32     `y3:"0x11"`
	Therm Thermometer `y3:"0x12"`
}

// Thermometer define struct for this example
type Thermometer struct {
	Temperature float32 `y3:"0x13"`
	Humidity    float32 `y3:"0x14"`
}

func main() {
	input := SourceData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: Thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	// encode interface to []byte
	codec := y3.NewCodec(0x20)
	inputBuf, _ := codec.Marshal(input)
	fmt.Printf("inputBuf=%#v\n", inputBuf)

	// define callback function to process the data being observed
	callback := func(v []byte) (interface{}, error) {
		return y3.ToFloat32(v)
	}

	// create the Observable interface
	source := y3.FromStream(bytes.NewReader(inputBuf))
	// subscribe observed to value
	consumer := source.Subscribe(0x11).OnObserve(callback)
	// checking data after it has been processed
	for c := range consumer {
		fmt.Printf("observed value=%v, type=%v\n", c, reflect.ValueOf(c).Kind())
	}
}
