package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/y3-codec-golang/internal/utils"
)

func main() {
	// encode
	data := []string{"abc", "def"}
	var node = y3.NewNodeSlicePacketEncoder(0x10)
	if out, ok := utils.ToStringSlice(data); ok {
		for _, v := range out {
			var item = y3.NewPrimitivePacketEncoder(0x00)
			item.SetStringValue(fmt.Sprintf("%v", v))
			node.AddPrimitivePacket(item)
		}
	}
	buf := node.Encode()
	// decode
	packet, _, _ := y3.DecodeNodePacket(buf)
	result := make([]string, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUTF8String()
		result = append(result, v)
	}
	fmt.Printf("result=%v", result)
}
