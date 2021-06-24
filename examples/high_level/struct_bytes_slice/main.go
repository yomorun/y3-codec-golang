package main

import (
	"bytes"
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding struct slice type with []byte
func main() {
	// Simulate source to generate and send data
	data := []ShakeData{
		{Topic: "shake", Payload: []byte{0x20, 0x21, 0x22}},
		{Topic: "shake2", Payload: []byte{0x30, 0x31, 0x32}},
	}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))

	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		var sl []ShakeData
		err := y3.ToObject(v, &sl)
		if err != nil {
			return nil, err
		}
		for _, s := range sl {
			fmt.Printf("encoded Payload: %#v\n", s.Payload)
		}
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
}

type ShakeData struct {
	Topic   string `y3:"0x11"`
	Payload []byte `y3:"0x12"`
}
