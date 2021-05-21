package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding bool type by using PrimitivePacket.
func main() {
	// encode
	var data bool = true
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetBoolValue(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToBool()
	fmt.Printf("val=%v", val)
}
