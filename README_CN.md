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

## Types

Y3实现了[YoMo Codec](https://github.com/yomorun/yomo-codec)协议，并支持如下Golang数据类型，同时提供了[Hight-Level](https://github.com/yomorun/y3-codec-golang/blob/master/explainer_CN.md#types)的封装。

<details>
  <summary>int32</summary>

```golang
// encode
var data int32 = 123
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetInt32Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToInt32()
fmt.Printf("val=%d", val)
```

</details>
<details>
  <summary>uint32</summary>

```golang
var data uint32 = 123
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetUInt32Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToUInt32()
fmt.Printf("val=%d", val)  
```

</details>
<details>
  <summary>int64</summary>

```golang
// encode
var data int64 = 123
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetInt64Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToInt64()
fmt.Printf("val=%d", val)  
```

</details>
<details>
  <summary>uint64</summary>

```golang
// encode
var data uint64 = 123
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetUInt64Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToUInt64()
fmt.Printf("val=%d", val) 
```

</details>
<details>
  <summary>float32</summary>

```golang
// encode
var data float32 = 1.23
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetFloat32Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToFloat32()
fmt.Printf("val=%f", val) 
```

</details>
<details>
  <summary>float64</summary>

```golang
// encode
var data float64 = 1.23
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetFloat64Value(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToFloat64()
fmt.Printf("val=%f", val)  
```

</details>
<details>
  <summary>bool</summary>

```golang
// encode
var data bool = true
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetBoolValue(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToBool()
fmt.Printf("val=%v", val)  
```

</details>
<details>
  <summary>string</summary>

```golang
// encode
var data string = "abc"
var prim = y3.NewPrimitivePacketEncoder(0x01)
prim.SetStringValue(data)
buf := prim.Encode()
// decode
res, _, _, _ := y3.DecodePrimitivePacket(buf)
val, _ := res.ToUTF8String()
fmt.Printf("val=%s", val) 
```

</details>
<details>
  <summary>int32 slice</summary>

```golang
// encode
data := []int32{123, 456}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToInt64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetInt32Value(int32(v.(int64)))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]int32, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToInt32()
	result = append(result, v)
}
fmt.Printf("result=%v", result) 
```

</details>
<details>
  <summary>uint32 slice</summary>

```golang
// encode
data := []uint32{123, 456}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToUInt64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetUInt32Value(uint32(v.(uint64)))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]uint32, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToUInt32()
	result = append(result, v)
}
fmt.Printf("result=%v", result) 
```

</details>
<details>
  <summary>int64 slice</summary>

```golang
// encode
data := []int64{123, 456}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToInt64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetInt64Value(v.(int64))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]int64, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToInt64()
	result = append(result, v)
}
fmt.Printf("result=%v", result) 
```

</details>
<details>
  <summary>uint64 slice</summary>

```golang
data := []uint64{123, 456}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToUInt64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetUInt64Value(v.(uint64))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]uint64, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToUInt64()
	result = append(result, v)
}
fmt.Printf("result=%v", result)  
```

</details>
<details>
  <summary>float32 slice</summary>

```golang
// encode
data := []float32{1.23, 4.56}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToUFloat64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetFloat32Value(float32(v.(float64)))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]float32, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToFloat32()
	result = append(result, v)
}
fmt.Printf("result=%v", result) 
```

</details>
<details>
  <summary>float64 slice</summary>

```golang
// encode
data := []float64{1.23, 4.56}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToUFloat64Slice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetFloat64Value(v.(float64))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]float64, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToFloat64()
	result = append(result, v)
}
fmt.Printf("result=%v", result)  
```

</details>
<details>
  <summary>bool slice</summary>

```golang
// encode
data := []bool{true, false}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToBoolSlice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetBoolValue(v.(bool))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]bool, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToBool()
	result = append(result, v)
}
fmt.Printf("result=%v", result)  
```

</details>
<details>
  <summary>string slice</summary>

```golang
// encode
data := []string{"abc", "def"}
var node = y3.NewNodeSlicePacketEncoder(0x10)
if out, ok := utils.ToStringSlice(data); ok {
	for _, v := range out {
		var item = y3.NewPrimitivePacketEncoder(0x00)
		item.SetStringValue(fmt.Sprintf("%v", v))
		node.AddPrimitivePacket(item)
	}
}
buf := node.Encode()
// decode
packet, _, _ := y3.DecodeNodePacket(buf)
result := make([]string, 0)
for _, p := range packet.PrimitivePackets {
	v, _ := p.ToUTF8String()
	result = append(result, v)
}
fmt.Printf("result=%v", result) 
```

</details>
<details>
  <summary>complex types</summary>

```golang
// encode
var node = y3.NewNodePacketEncoder(0x01)
node.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
	var prim1 = y3.NewPrimitivePacketEncoder(0x10)
	prim1.SetFloat32Value(40.5)
	return prim1
}())
node.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
	var prim1 = y3.NewPrimitivePacketEncoder(0x11)
	prim1.SetInt64Value(time.Now().Unix())
	return prim1
}())
buf := node.Encode()
// decode
res, _, _ := y3.DecodeNodePacket(buf)
for _, v := range res.PrimitivePackets {
	if v.SeqID() == 0x10 {
		fmt.Printf("0x10=%f\n", func() float32 {
			val, _ := v.ToFloat32()
			return val
		}())
	}
	if v.SeqID() == 0x11 {
		fmt.Printf("0x11=%d\n", func() int64 {
			val, _ := v.ToInt64()
			return val
		}())
	}
}  
```

</details>
<details>
  <summary>complex slice types</summary>

```golang
// encode
var node = y3.NewNodeSlicePacketEncoder(0x01)
for i := 0; i < 2; i++ {
	item := y3.NewNodePacketEncoder(0x00)
	item.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
		var prim1 = y3.NewPrimitivePacketEncoder(0x10)
		prim1.SetFloat32Value(40.5)
		return prim1
	}())
	item.AddPrimitivePacket(func() *y3.PrimitivePacketEncoder {
		var prim1 = y3.NewPrimitivePacketEncoder(0x11)
		prim1.SetInt64Value(time.Now().Unix())
		return prim1
	}())
	node.AddNodePacket(item)
}
buf := node.Encode()
// decode
res, _, _ := y3.DecodeNodePacket(buf)
for _, v := range res.NodePackets {
	if res.SeqID() != 0x01 {
		continue
	}
	for _, vv := range v.PrimitivePackets {
		if vv.SeqID() == 0x10 {
			fmt.Printf("0x10=%f\n", func() float32 {
				val, _ := vv.ToFloat32()
				return val
			}())
		}
		if vv.SeqID() == 0x11 {
			fmt.Printf("0x11=%d\n", func() int64 {
				val, _ := vv.ToInt64()
				return val
			}())
		}
	}
} 
```

</details>


## Contributors

[//]: contributor-faces

<a href="https://github.com/figroc"><img src="https://avatars1.githubusercontent.com/u/2026460?v=3" title="figroc" width="80" height="80"></a>

[//]: contributor-faces

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fyomorun%2Fy3-codec-golang?ref=badge_large)

