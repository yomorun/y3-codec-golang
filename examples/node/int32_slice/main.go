package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// Example of encoding and decoding int32 slice type by using NodePacket.
func main() {
	// encode
	data := []int32{123, 456}
	var node = y3.NewNodeSlicePacketEncoder(0x10)
	if out, ok := utils.ToInt64Slice(data); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(0x00)
			item.SetInt32Value(int32(v.(int64)))
			node.AddPrimitivePacket(item)
		}
	}
	buf := node.Encode()
	// decode
	packet, _, _ := y3.DecodeNodePacket(buf)
	result := make([]int32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt32()
		result = append(result, v)
	}
	fmt.Printf("result=%v", result)
}
