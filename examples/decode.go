package main

import (
	"fmt"

	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	fmt.Println("hello YoMo Codec golang implementation: Y3")
	encodePacket()
	// parseNodePacket()
	// parseStringPrimitivePacket()
}

type bar struct {
	Name string
}

type foo struct {
	ID int
	*bar
}

func encodePacket() {
	// We will encode JSON-like object `obj`:
	// 0x81: {
	//   0x02: -1,
	//   0x83 : {
	//     0x04: "C",
	//   },
	// }
	// to
	// [0x81, 0x08, 0x02, 0x01, 0x7F, 0x83, 0x03, 0x04, 0x01, 0x43]
	var obj = &foo{ID: -1, bar: &bar{Name: "C"}}

	// 0x81 - node
	var yFoo = y3.NewNodePacketEncoder(0x01)

	// 0x02 - ID=1
	var yp1 = y3.NewPrimitivePacketEncoder(0x02)
	yp1.SetInt32Value(-1)
	yFoo.AddPrimitivePacket(yp1)

	// 0x83 - &bar{}
	var yBar = y3.NewNodePacketEncoder(0x03)

	// 0x04 - Name: "C"
	var yp2 = y3.NewPrimitivePacketEncoder(0x04)
	yp2.SetStringValue("C")
	yBar.AddPrimitivePacket(yp2)

	yFoo.AddNodePacket(yBar)

	fmt.Println(obj)
	fmt.Printf("res=%#v", yFoo.Encode())
}

func parseNodePacket() {
	fmt.Println(">> Parsing [0x84, 0x0A, 0x0A, 0x01, 0x7F, 0x0B, 0x05, 0x43, 0x45, 0x4C, 0x4C, 0x41] EQUALS JSON= 0x84: { 0x0A: -1, 0x0B: 'CELLA' }")
	buf := []byte{0x84, 0x0A, 0x0A, 0x01, 0x7F, 0x0B, 0x05, 0x43, 0x45, 0x4C, 0x4C, 0x41}
	res, _, err := y3.DecodeNodePacket(buf)
	v1 := res.PrimitivePackets[0]

	p1, err := v1.ToInt32()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Tag Key=[%#X.%#X], Value=%v\n", res.SeqID(), v1.SeqID(), p1)

	v2 := res.PrimitivePackets[1]

	p2, err := v2.ToUTF8String()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X.%#X], Value=%v\n", res.SeqID(), v2.SeqID(), p2)
}

func parseStringPrimitivePacket() {
	fmt.Println(">> Parsing [0x0A, 0x01, 0x7F] EQUALS JSON= { 0x0A: 127 } ")
	buf := []byte{0x0A, 0x01, 0x7F}
	res, _, err := y3.DecodePrimitivePacket(buf)
	v1, err := res.ToUInt32()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X], Value=%v\n", res.SeqID(), v1)
}
