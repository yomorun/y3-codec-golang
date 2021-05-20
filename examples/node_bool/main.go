package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func main() {
	// encode
	data := []bool{true, false}
	var node = y3.NewNodeSlicePacketEncoder(0x10)
	if out, ok := utils.ToBoolSlice(data); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(0x00)
			item.SetBoolValue(v.(bool))
			node.AddPrimitivePacket(item)
		}
	}
	buf := node.Encode()
	// decode
	packet, _, _ := y3.DecodeNodePacket(buf)
	result := make([]bool, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToBool()
		result = append(result, v)
	}
	fmt.Printf("result=%v", result)
}
