package main

import (
	"fmt"
	"time"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding complex slice types by using NodePacket.
func main() {
	// encode
	var node = y3.NewNodeSlicePacketEncoder(0x01)
	for i := 0; i < 2; i++ {
		item := y3.NewNodePacketEncoder(0x00)
		item.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
			var prim1 = y3.NewPrimitivePacketEncoder(0x10)
			prim1.SetFloat32Value(40.5)
			return prim1
		}())
		item.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
			var prim1 = y3.NewPrimitivePacketEncoder(0x11)
			prim1.SetInt64Value(time.Now().Unix())
			return prim1
		}())
		node.AddNodePacket(item)
	}
	buf := node.Encode()
	// decode
	res, _, _ := y3.DecodeNodePacket(buf)
	for _, v := range res.NodePackets {
		if res.SeqID() != 0x01 {
			continue
		}
		for _, vv := range v.PrimitivePackets {
			if vv.SeqID() == 0x10 {
				fmt.Printf("0x10=%f\n", func() float32 {
					val, _ := vv.ToFloat32()
					return val
				}())
			}
			if vv.SeqID() == 0x11 {
				fmt.Printf("0x11=%d\n", func() int64 {
					val, _ := vv.ToInt64()
					return val
				}())
			}
		}
	}
}
