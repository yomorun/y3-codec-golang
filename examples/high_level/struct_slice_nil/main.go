package main

import (
	"bytes"
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

func main() {
	// Simulate source to generate and send data, value is nil.
	var data []NoiseData
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))

	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		var sl []NoiseData
		err := y3.ToObject(v, &sl)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded slice len: %v\n", len(sl))
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
}

type NoiseData struct {
	Noise float32 `y3:"0x11"`
	Other []int64 `y3:"0x12"`
}
