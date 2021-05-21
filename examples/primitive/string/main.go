package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding string type by using PrimitivePacket.
func main() {
	// encode
	var data string = "abc"
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetStringValue(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToUTF8String()
	fmt.Printf("val=%s", val)
}
