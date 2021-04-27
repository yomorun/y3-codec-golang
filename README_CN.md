> 📚 VERSION: draft-02
>
> ⛳️ STATE: WIP
>
> 🇨🇳 [简体中文](https://github.com/yomorun/y3-codec-golang/blob/master/explainer_CN.md)  🇬🇧 [English](https://github.com/yomorun/y3-codec-golang/blob/master/README.md)

# Y3
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang?ref=badge_shield)

Y3是一种[YoMo Codec](https://github.com/yomorun/yomo-codec)的Golang实现，它描述了个快速和低CPU损耗的编解码器，专注于边缘计算和流处理。

查看 [explainer](https://github.com/yomorun/y3-codec-golang/blob/master/explainer_CN.md) 获取更多信息，了解更多与[YoMo](https://github.com/yomorun/yomo)组合的方式。

## Y3 Codec

See [Y3 Codec SPEC](https://github.com/yomorun/y3-codec)

## Test

```
go test ./...
```

## Use

```
go get -u github.com/yomorun/y3-codec-golang
```

## Examples

### 编码例子

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
)

func main() {
	// if we want to repesent `var obj = &foo{ID: -1, bar: &bar{Name: "C"}}` 
	// in Y3-Codec:

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

### 解码例子 1:  解码一个原始数据包

```go
package main

import (
	"fmt"
	y3 "github.com/yomorun/y3-codec-golang"
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

### 解码例子 2: 解码一个节点数据包

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

## Performance

- 与JSON的性能比较: [yomo-y3-stress-testing](https://github.com/10cella/yomo-y3-stress-testing)
- 与ProtoBuffer的性能比较：[y3-protobuf-testing](https://github.com/yomorun/y3-protobuf-testing)

## Contributors

[//]: contributor-faces

<a href="https://github.com/figroc"><img src="https://avatars1.githubusercontent.com/u/2026460?v=3" title="figroc" width="80" height="80"></a>

[//]: contributor-faces

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang?ref=badge_large)

