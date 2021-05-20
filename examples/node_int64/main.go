package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func main() {
	// encode
	data := []int64{123, 456}
	var node = y3.NewNodeSlicePacketEncoder(0x10)
	if out, ok := utils.ToInt64Slice(data); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(0x00)
			item.SetInt64Value(v.(int64))
			node.AddPrimitivePacket(item)
		}
	}
	buf := node.Encode()
	// decode
	packet, _, _ := y3.DecodeNodePacket(buf)
	result := make([]int64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt64()
		result = append(result, v)
	}
	fmt.Printf("result=%v", result)
}
