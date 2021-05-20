package main

import (
	"fmt"
	"time"

	"github.com/yomorun/y3-codec-golang"
)

func main() {
	// encode
	var node = y3.NewNodePacketEncoder(0x01)
	node.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
		var prim1 = y3.NewPrimitivePacketEncoder(0x10)
		prim1.SetFloat32Value(40.5)
		return prim1
	}())
	node.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
		var prim1 = y3.NewPrimitivePacketEncoder(0x11)
		prim1.SetInt64Value(time.Now().Unix())
		return prim1
	}())
	buf := node.Encode()
	// decode
	res, _, _ := y3.DecodeNodePacket(buf)
	for _, v := range res.PrimitivePackets {
		if v.SeqID() == 0x10 {
			fmt.Printf("0x10=%f\n", func() float32 {
				val, _ := v.ToFloat32()
				return val
			}())
		}
		if v.SeqID() == 0x11 {
			fmt.Printf("0x11=%d\n", func() int64 {
				val, _ := v.ToInt64()
				return val
			}())
		}
	}
}
