# Y3Codec Explainer

## Problem and Motivation

- There isn't an efficient codec under persistent connection in the market, for example Protobuf is aimed at the scene of offline data, while users are facing the real-time data processing, and this kind of data is generated at high-frequency and the changing of data structure is low-frequency.
- In the decoding stage, other codecs are full decode, we have to get a full data packet before we can start the decoding process, which is low efficiency in real-time data processing.

## Goals

- Faster than real-time.
- Provide efficient `decoding` for real-time stream processing under persistent connection. Users don't have to wait to get the full packet before they can start `decoding`.
- Users only need to 'observe key' to get the data they concern.

## Key use-cases

- Low-latency sensitive applications.
- Real-time stream processing under persistent connection.

## Proposed solutions

`Y3`'s scenario is to deal with real-time stream processing under persistent connection, so the user gives the raw stream to `Y3`, and then tells `Y3` to observe the key. `Y3` starts the parsing operation after taking over the raw stream. When the key is observed, `Y3` will decode the value in specified type and trigger callback function.

Core interfaces include:

- **Marshal** serializes the user's data according to the `Y3`'s encoding rules.
- **Subscribe** observes the `key` which the user specified.
- **OnObserve** triggers the callback function while the is observed by `Y3`.

## Examples

### 1. The data source is a batch of JSON (including concerned and unconcerned data). These data need to be encoded by 'Y3' and transported to the receiver by streaming, such as `YoMo flow`. The receiver observes the concerned data and processes it

#### Encode data

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
   codec := y3.NewCodec(0x20)
   inputBuf, err := codec.Marshal(input)
   fmt.Printf("inputBuf=%v, err=%v\n", packetutils.FormatBytes(inputBuf), err)
}
```

#### Observe the specified key and decode

```go
func main() {
   var onObserve = func(v []byte) {
      f, err := y3.ToFloat32(v)
      fmt.Printf("observed v=%v\n", f)
   }
   codec = y3.FromStream(xx)
   codec.Subscribe(0x11).OnObserve(onObserve)
}
```
