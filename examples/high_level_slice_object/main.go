package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/yomorun/y3-codec-golang"
)

/**
Simulate how to encode and decode the slice of struct type

Decode by method:
y3.ToObject
*/
func main() {
	// Simulate source to generate and send data
	data := []NoiseData{
		{Noise: 40, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
		{Noise: 50, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
	}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))

	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		var sl []NoiseData
		err := y3.ToObject(v, &sl)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
}

type NoiseData struct {
	Noise float32 `y3:"0x11"`
	Time  int64   `y3:"0x12"`
	From  string  `y3:"0x13"`
}
