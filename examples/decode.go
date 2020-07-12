package main

import (
	"fmt"

	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	fmt.Printf("hello YoMo Codec golang implementation: Y3\n")
	parseNodePacket()
	parseStringPrimitivePacket()
}

func parseNodePacket() {
	fmt.Println(">> Parsing [0x84, 0x08, 0x01, 0x04, 0x01, 0x01]")
	buf := []byte{0x84, 0x08, 0x01, 0x04, 0x01, 0x01}
	res, _, err := y3.DecodeNodePacket(buf)
	v1 := res.PrimitivePackets[0]

	p1, err := v1.ToInt64()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X.%#X], Value=%v\n", res.Tag.SeqID(), v1.Tag, p1)
}

func parseStringPrimitivePacket() {
	fmt.Println(">> Parsing [0x0B, 0x0C, 0x00, 0x43, 0x45, 0x4C, 0x4C, 0x41]")
	buf := []byte{0x0B, 0x0C, 0x00, 0x43, 0x45, 0x4C, 0x4C, 0x41}
	res, _, err := y3.DecodePrimitivePacket(buf)
	v1, err := res.ToUTF8String()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X], Value=%v", res.Tag, v1)
}
