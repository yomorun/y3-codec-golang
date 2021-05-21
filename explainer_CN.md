# Y3Codec Explainer

## Problem and Motivation

- 市面上在长连接下缺乏正确的 codec，Protobuf 也是针对离线数据的场景，而用户面对的是未来实时的数据处理，而这类数据的特点就是高频产生，数据结构的变化频次小。
- 在解码阶段，其他的 codec 都是 Fully Decode，即一定要拿到一个完整的数据包，才能开始反序列化过程，这在实时数据处理环节是大忌。

## Goals

- Faster than Real-time
- 针对长连接下的实时数据流处理，提供高效的`解码`。用户无需等待拿到完整的数据包之后，才能开始`解码`。
- 用户只需 `observe key` 即可拿到想要的数据。

## Key use-cases

- 对延迟敏感的应用程序。
- 长连接下的实时数据流处理。

## Proposed solutions

`Y3` 的场景是应对长连接下的实时数据流处理，所以用户把 raw stream 交给 `Y3`，然后告诉 `Y3` 要 observe 的 key，`Y3` 在接管 raw stream 后开始 parsing 的操作，发现了 key 后，开始将其对应的 value 以用户指定的数据类型做反序列化，再将其作为参数，调用用户指定的回调函数（event-driven method）。

`Y3` 通过把对象数据描述成一组 `TLV` 结构，在数据包解码时，可以在解码过程中更早的了解到当前的 T 是否为所受监听的 key，从而判断是否直接跳到下一组TLV结构，而并不需要对非受监听的数据包进行多余的解码操作，从而提升了解码的效率和资源利用率。

主要接口包括：

- **Marshal** 按照 `Y3` 的编码规则序列化用户的数据。
- **Subscribe** 监听用户想 observe 的 `key`。
- **OnObserve** `Y3` 发现了key后，调用用户指定的回调函数进行解码。

## Examples

### 1.数据源为一批类拟 JSON 的层级数据(其中包含了关心和不关心的数据)，需要把这些数据转换成 `Y3` 编码通过流式传输给接收方，比如 [yomo-flow](https://yomo.run/flow)。接收方监听关心的数据并进行业务处理。

#### 编码数据

```go
type SourceData struct {
   Name  string      `y3:"0x10"`
   Noise float32     `y3:"0x11"`
   Therm Thermometer `y3:"0x12"`
}

type Thermometer struct {
   Temperature float32 `y3:"0x13"`
   Humidity    float32 `y3:"0x14"`
}

func main() {
	input := SourceData{
		Name:  "yomo",
		Noise: float32(456),
		Therm: Thermometer{Temperature: float32(30), Humidity: float32(40)},
	}
  // 把对象编码为符合Y3-Codec格式的数据
	codec := y3.NewCodec(0x20)
	inputBuf, _ := codec.Marshal(input)
	fmt.Printf("inputBuf=%#v\n", inputBuf)
}
```

#### 解码并监听一个值

```go
func main() {
	// 定义回调函数用于处理被监听的数据
	callback := func(v []byte) (interface{}, error) {
		return y3.ToFloat32(v)
	}
	// 创建Observable接口
	source := y3.FromStream(bytes.NewReader(inputBuf))
	// 订阅被监听的Key，并设置回调函数
	consumer := source.Subscribe(0x11).OnObserve(callback)
	// 检查被处理后的数据
	for c := range consumer {
		fmt.Printf("observed value=%v, type=%v\n", c, reflect.ValueOf(c).Kind())
	}
}
```



## Types

Y3提供High Level的封装，用于支持像YoMo这样的框架。

统一的编码方法: `y3.NewCodec(observe byte).Marshal(input interface{})`

| 类型          | 解码方法             |
| ------------- | -------------------- |
| struct        | y3.ToObject          |
| struct slice  | y3.ToObject          |
| int32         | y3.ToInt32           |
| int32 slice   | y3.ToInt32Slice      |
| uint32        | y3.ToUInt32          |
| uint32 slice  | y3.ToUInt32Slice     |
| int64         | y3.ToInt64           |
| int64 slice   | y3.ToInt64Slice      |
| uint64        | y3.ToUInt64          |
| uint64 slice  | y3.ToUInt64Slice     |
| float32       | y3.ToFloat32         |
| float32 slice | y3.ToFloat32Slice    |
| float64       | y3.ToFloat64         |
| float64 slice | y3.ToFloat64Slice    |
| bool          | y3.ToBool            |
| bool slice    | y3.ToBoolSlice       |
| string        | y3.ToUTF8String      |
| string slice  | y3.ToUTF8StringSlice |

<details>
  <summary>struct</summary>
  <pre class="go" style="background-color: aliceblue">
  func main() {
    // Simulate source to generate and send data
    data := <b>NoiseData</b>{Noise: 40, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"}
    sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
    source := y3.FromStream(bytes.NewReader(sendingBuf))
    // Simulate flow listening and decoding data
    var decode = func(v []byte) (interface{}, error) {
        var obj NoiseData
        err := y3.<b>ToObject</b>(v, <b>&obj</b>)
        if err != nil {
            return nil, err
        }
        fmt.Printf("encoded data: %v\n", obj)
        return obj, nil
    }
    consumer := source.Subscribe(0x10).OnObserve(decode)
    for range consumer {
    }
  }
  type <b>NoiseData</b> struct {
      Noise float32 `y3:"0x11"`
      Time  int64   `y3:"0x12"`
      From  string  `y3:"0x13"`
  }
  </pre>
</details>
<details>
  <summary>struct slice</summary>
  <pre class="go" style="background-color: aliceblue">
    func main() {
      // Simulate source to generate and send data
      data := <b>[]NoiseData</b>{
          {Noise: 40, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
          {Noise: 50, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
      }
      sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
      source := y3.FromStream(bytes.NewReader(sendingBuf))
      // Simulate flow listening and decoding data
      var decode = func(v []byte) (interface{}, error) {
          var sl []NoiseData
          err := y3.<b>ToObject</b>(v, <b>&sl</b>)
          if err != nil {
              return nil, err
          }
          fmt.Printf("encoded data: %v\n", sl)
          return sl, nil
      }
      consumer := source.Subscribe(0x10).OnObserve(decode)
      for range consumer {
      }
    }
    type <b>NoiseData</b> struct {
        Noise float32 `y3:"0x11"`
        Time  int64   `y3:"0x12"`
        From  string  `y3:"0x13"`
    }
  </pre>
</details>
<details>
  <summary>int32</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>int32</b> = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToInt32</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>int32 slice</summary>
  <pre class="go" style="background-color: aliceblue">
    // Simulate source to generate and send data
    data := []<b>int32</b>{123, 456}
    sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
    source := y3.FromStream(bytes.NewReader(sendingBuf))
    // Simulate flow listening and decoding data
    var decode = func(v []byte) (interface{}, error) {
        sl, err := y3.<b>ToInt32Slice</b>(v)
        if err != nil {
            return nil, err
        }
        fmt.Printf("encoded data: %v\n", sl)
        return sl, nil
    }
    consumer := source.Subscribe(0x10).OnObserve(decode)
    for range consumer {
    }
  </pre>
</details>
<details>
  <summary>uint32</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>uint32</b> = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToUInt32</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>uint32 slice</summary>
  <pre class="go" style="background-color: aliceblue">
  // Simulate source to generate and send data
  data := []<b>uint32</b>{123, 456}
  sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
  source := y3.FromStream(bytes.NewReader(sendingBuf))
  // Simulate flow listening and decoding data
  var decode = func(v []byte) (interface{}, error) {
      sl, err := y3.<b>ToUInt32Slice</b>(v)
      if err != nil {
          return nil, err
      }
      fmt.Printf("encoded data: %v\n", sl)
      return sl, nil
  }
  consumer := source.Subscribe(0x10).OnObserve(decode)
  for range consumer {
  }
  </pre>
</details>
<details>
  <summary>int64</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>int64</b> = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToInt64</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>int64 slice</summary>
  <pre class="go" style="background-color: aliceblue">
  // Simulate source to generate and send data
  data := []<b>int64</b>{123, 456}
  sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
  source := y3.FromStream(bytes.NewReader(sendingBuf))
  // Simulate flow listening and decoding data
  var decode = func(v []byte) (interface{}, error) {
      sl, err := y3.<b>ToInt64Slice</b>(v)
      if err != nil {
          return nil, err
      }
      fmt.Printf("encoded data: %v\n", sl)
      return sl, nil
  }
  consumer := source.Subscribe(0x10).OnObserve(decode)
  for range consumer {
  }
  </pre>
</details>
<details>
  <summary>uint64</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>uint64</b> = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToUInt64</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>uint64 slice</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := []<b>uint64</b>{123, 456}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToUInt64Slice</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>float32</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>float32</b> = 1.23
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToFloat32</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>float32 slice</summary>
  <pre class="go" style="background-color: aliceblue">
  // Simulate source to generate and send data
	data := []<b>float32</b>{1.23, 4.56}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToFloat32Slice</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>float64</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	var data <b>float64</b> = 1.23
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToFloat64</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>float64 slice</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := []<b>float64</b>{1.23, 4.56}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToFloat64Slice</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>bool</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := true
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToBool</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>bool slice</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := []<b>bool</b>{true, false}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToBoolSlice</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>string</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := "abc"
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToUTF8String</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>
<details>
  <summary>string slice</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := []<b>string</b>{"a", "b"}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.<b>ToUTF8StringSlice</b>(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  </pre>
</details>

更多例子请看： `/examples/`

## Attention

被监听的key有规定的使用范围：
- 用户自定义的被监听key范围：`0x10 ~ 0x3f`
- 系统保留： `0x01 ~ 0x0f`

