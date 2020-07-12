# yomo-codec-golang

still *Working In Progress*

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
```