package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding uint32 type by using PrimitivePacket.
func main() {
	// encode
	var data uint32 = 123
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetUInt32Value(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToUInt32()
	fmt.Printf("val=%d", val)
}
