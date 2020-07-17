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

```go
package main

import (
	"fmt"

	y3 "github.com/yomorun/yomo-codec-golang"
)

func main() {
	fmt.Println("hello YoMo Codec golang implementation: Y3")
	fmt.Println(">> Parsing [0x81, 0x03, 0x01, 0x01, 0x7F]")
	buf := []byte{0x81, 0x03, 0x01, 0x01, 0x7F}
	res, _, err := y3.DecodeNodePacket(buf)
	v1 := res.PrimitivePackets[0]

	p1, err := v1.ToInt64()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X], Value=%v\n", res.SeqID(), p1)
}
```

More examples in `/examples/`

## YoMo Codec

See [SPEC](https://github.com/yomorun/yomo-codec)