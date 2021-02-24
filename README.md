> 📚 VERSION: draft-02
>
> ⛳️ STATE: WIP
>
> 🇨🇳 [简体中文](https://gitee.com/yomorun/y3-codec-golang/blob/master/README_CN.md)  🇬🇧 [English](https://github.com/yomorun/y3-codec-golang/blob/master/README.md)

# Y3
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang?ref=badge_shield)

Y3 is the golang implementation of [YoMo Codec](https://github.com/yomorun/yomo-codec), which describe a fast and low CPU binding data encoder/decoder focus on edge computing and streaming processing.

See the [explainer](explainer.md) for more info.

## YoMo Codec

See [YoMo Codec SPEC](https://github.com/yomorun/yomo-codec)

## Test

`make test`

## Use 

`go get -u github.com/yomorun/y3-codec-golang`

## Examples

### Encode examples

Encode Key-Value type `{ID: 3}` to Y3: 

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
)

func main() {
	// Key:ID = Tag:0x01
	buffer, _ = y3.EncodeInt(0x01, 3)
	
	// get whole buf
	fmt.Printf("res=%#v", buffer) // res=[]byte{0x01, 0x01, 0x03}
}
```

if we want to repesent JSON `foo: {ID: -1, Name: "C"}` in Y3:

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
)

func main() {
	// Key:ID 0x01 -> -1
	var id, _ = y3.NewPacket(0x01)
	id.SetInt32(-1)
	
	// Key:Name 0x02 -> "C"
	var name, _ = y3.NewPacket(0x02)
	name.SetString("C")

	// parent node
	var foo, _ = y3.NewPacket(0x10)
	foo.AddNode(id)
	foo.AddNode(name)
	
	// get whole buf
	fmt.Printf("res=%#v", foo.Encode()) // res=[]byte{0x10, 0x06, 0x01, 0x01, 0x7F, 0x02, 0x01, 0x43}
}
```

### Decode examples 1: decode a primitive packet

Decode `[0x0A, 0x01, 0x7F]` as Int type

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
)

func main() {
	fmt.Println(">> Parsing [0x0A, 0x01, 0x7F], which like Key-Value format = 0x0A: 127")
	buf := []byte{0x0A, 0x01, 0x7F}
	res, _, err := y3.DecodePacket(buf)
	val, err := res.ToInt32()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Tag Key=[%#X], Value=%v\n", res.Tag, val)
}
```

### Decode examples 2: decode a node packet

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
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

## Contributors

[//]: contributor-faces

<a href="https://github.com/figroc"><img src="https://avatars1.githubusercontent.com/u/2026460?v=3" title="figroc" width="80" height="80"></a>

[//]: contributor-faces

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang?ref=badge_large)

## CLA

[Sign CLA](https://cla-assistant.io/yomorun/y3-codec-golang)

[![CLA assistant](https://cla-assistant.io/readme/badge/yomorun/y3-codec-golang)](https://cla-assistant.io/yomorun/y3-codec-golang)