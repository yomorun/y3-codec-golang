> 📚 VERSION: draft-01
>
> ⛳️ STATE: WIP

# Y3

Golang implementation of [YoMo Codec](https://github.com/yomorun/yomo-codec)

## Test

`go test ./...`

## Use 

`go get -u github.com/yomorun/yomo-codec-golang`

## Examples

### Encode examples

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	// if we want to repesent `var obj = &foo{ID: -1, bar: &bar{Name: "C"}}` 
	// in YoMo-Codec:

	// 0x81 -> node
	var foo = y3.NewNodePacketEncoder(0x01)

	// 0x02 -> foo.ID=-11
	var yp1 = y3.NewPrimitivePacketEncoder(0x02)
	yp1.SetInt32Value(-1)
	foo.AddPrimitivePacket(yp1)

	// 0x83 -> &bar{}
	var bar = y3.NewNodePacketEncoder(0x03)

	// 0x04 -> bar.Name="C"
	var yp2 = y3.NewPrimitivePacketEncoder(0x04)
	yp2.SetStringValue("C")
	bar.AddPrimitivePacket(yp2)
	
	// -> foo.bar=&bar
	foo.AddNodePacket(bar)

	fmt.Printf("res=%#v", foo.Encode()) // res=[]byte{0x81, 0x08, 0x02, 0x01, 0x7F, 0x83, 0x03, 0x04, 0x01, 0x43}
}
```

### Decode examples 1: decode a primitive packet

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	fmt.Println(">> Parsing [0x0A, 0x01, 0x7F], which like Key-Value format = 0x0A: 127")
	buf := []byte{0x0A, 0x01, 0x7F}
	res, _, err := y3.DecodePrimitivePacket(buf)
	v1, err := res.ToUInt32()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X], Value=%v\n", res.SeqID(), v1)
}
```

### Decode examples 2: decode a node packet

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	fmt.Println(">> Parsing [0x84, 0x06, 0x0A, 0x01, 0x7F, 0x0B, 0x01, 0x43] EQUALS JSON= 0x84: { 0x0A: -1, 0x0B: 'C' }")
	buf := []byte{0x84, 0x06, 0x0A, 0x01, 0x7F, 0x0B, 0x01, 0x43}
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
```

More examples in `/examples/`

## YoMo Codec

See [SPEC](https://github.com/yomorun/yomo-codec)
