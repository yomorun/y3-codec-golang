package main

import (
	"bytes"
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding bool slice type, value is nil.
func main() {
	// Simulate source to generate and send data
	var data []bool
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToBoolSlice(v)
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
