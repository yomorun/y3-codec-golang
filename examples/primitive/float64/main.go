package main

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang"
)

// Example of encoding and decoding float64 type by using PrimitivePacket.
func main() {
	// encode
	var data float64 = 1.23
	var prim = y3.NewPrimitivePacketEncoder(0x01)
	prim.SetFloat64Value(data)
	buf := prim.Encode()
	// decode
	res, _, _, _ := y3.DecodePrimitivePacket(buf)
	val, _ := res.ToFloat64()
	fmt.Printf("val=%f", val)
}
