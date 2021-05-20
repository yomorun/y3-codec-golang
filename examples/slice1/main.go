package main

import (
	"bytes"
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

/**
Simulate how to encode and decode the slice of int32 type

The supported types of slice:
[]int32,[]uint32,[]int64,[]uint64,[]float32,[]float64,[]bool,[]string

Use the following method for decoding：
y3.ToInt32Slice
y3.ToUInt32Slice
y3.ToInt64Slice
y3.ToUInt64Slice
y3.ToFloat32Slice
y3.ToFloat64Slice
y3.ToBoolSlice
y3.ToUTF8StringSlice
*/
func main() {
	// Simulate source to generate and send data
	data := []int32{123, 456}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToInt32Slice(v)
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
