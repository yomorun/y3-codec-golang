package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding uint64 type by using PrimitivePacket.
func main() {
	// encode
	var data uint64 = 123
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetUInt64Value(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToUInt64()
	fmt.Printf("val=%d", val)
}
