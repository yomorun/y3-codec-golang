package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

func main() {
	// encode
	var data int32 = 123
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetInt32Value(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToInt32()
	fmt.Printf("val=%d", val)
}
