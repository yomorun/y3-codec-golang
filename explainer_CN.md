# Y3Codec Explainer

## Problem and Motivation

- 市面上在长连接下缺乏正确的 codec，Protobuf 也是针对离线数据的场景，而用户面对的是未来实时的数据处理，而这类数据的特点就是高频产生，数据结构的变化频次小。
- 在解码阶段，其他的 codec 都是 Fully Decode，即一定要拿到一个完整的数据包，才能开始反序列化过程，这在实时数据处理环节是大忌。

## Goals

- 针对长连接下的实时数据流处理，提供高效的`解码`。用户无需等待拿到完整的数据包之后，才能开始`解码`。
- 用户只需 `observe key` 即可拿到想要的数据。

## Key use-cases

- 对延迟敏感的应用程序。
- 长连接下的实时数据流处理。

## Proposed solutions

`Y3` 的场景是应对长连接下的实时数据流处理，所以用户把 raw stream 交给 `Y3`，然后告诉 `Y3` 要 observe 的key，`Y3` 在接管 raw stream 后开始 parsing的 操作，发现了key后，开始将其对应的 value 以用户指定的数据类型做反序列化，再将其作为参数，调用用户指定的回调函数（event-driven method）。

主要接口包括：

- **Marshal** 按照 `Y3` 的编码规则序列化用户的数据。
- **Subscribe** 监听用户想 observe 的 `key`。
- **OnObserve** `Y3` 发现了key后，调用用户指定的回调函数。
- **OnUnobserved** `Y3` 未监听到用户指定的 key，调用指定的回调函数。

## Examples

### 数据源为一批类拟 JSON 的层级数据(其中包含了关心和不关心的数据)，需要把这些数据转换成 `Y3` 编码通过流式传输给接收方，比如 `yomo-flow`。接收方监听关心的数据并进行业务处理。

#### 编码数据

```go
type SourceData struct {
   Name  string      `yomo:"0x10"`
   Noise float32     `yomo:"0x11"`
   Therm Thermometer `yomo:"0x12"`
}

type Thermometer struct {
   Temperature float32 `yomo:"0x13"`
   Humidity    float32 `yomo:"0x14"`
}

func main() {
   input := SourceData{
      Name:  "yomo",
      Noise: float32(456),
      Therm: Thermometer{Temperature: float32(30), Humidity: float32(40)},
   }
   y3Codec := y3.NewY3Codec(0x20)
   inputBuf, err := y3Codec.Marshal(input)
   fmt.Printf("inputBuf=%v, err=%v\n", packetutils.FormatBytes(inputBuf), err)
}
```

#### 解码并监听一个值

```go
func main() {
   var onObserve = func(v []byte) {
      f, err := y3.ToFloat64(v)
      fmt.Printf("observed v=%v\n", f)
   }
   var onUnobserved = func(v []byte) {
      fmt.Printf("unobserved v=%v\n", v)
   }
   codec = y3.FromStream(xx)
   codec.Subscribe(0x10).OnObserve(onObserve).OnUnobserved(onUnobserved)
}
```
