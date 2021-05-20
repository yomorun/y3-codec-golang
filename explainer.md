# Y3Codec Explainer

## Problem and Motivation

- There isn't an efficient codec under persistent connection in the market, for example Protocol Buffers is aimed at the scene of offline data, while users are facing the real-time data processing, and this kind of data is generated at high-frequency and the changing of data structure is low-frequency.
- In the decoding stage, other codecs are full decode, the users have to get a full data packet before they can start the decoding process, which is low efficiency in real-time data processing.

## Goals

- Faster than real-time.
- Provide efficient `decoding` for real-time stream processing under persistent connection. Users don't have to wait to get the full packet before they can start `decoding`.
- Users only need to `observe key` to get the data they concern.

## Key use-cases

- Low-latency sensitive applications.
- Real-time stream processing under persistent connection.

## Proposed solutions

`Y3`'s scenario is to deal with real-time stream processing under persistent connection, so the user gives the raw stream to `Y3`, and then tells `Y3` to observe the key. `Y3` starts the parsing operation after taking over the raw stream. When the key is observed, `Y3` will decode the value in specified type and trigger a callback function.

`Y3` describes the object as a set of `TLV` structures. When decoding packets, it can know earlier whether the current T is the observed key, and determine whether to jump to the next set of TLV structures directly, without unnecessary decoding operations on unobserved packets, so `Y3` can improve the decoding efficiency and resource utilization.

Core interfaces include:

- **Marshal** serializes the user's data according to the `Y3`'s encoding rules.
- **Subscribe** observes the `key` which is specified by user.
- **OnObserve** triggers the callback function while the key is observed by `Y3`.

## Examples

### 1. The data source is a batch of JSON (including concerned and unconcerned data). These data need to be encoded by 'Y3' and transported to the receiver by streaming, such as [yomo-flow](https://yomo.run/flow). The receiver observes the concerned data and processes it

#### Encode data

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
  // encode to the data in Y3-Codec format
	codec := y3.NewCodec(0x20)
	inputBuf, _ := codec.Marshal(input)
	fmt.Printf("inputBuf=%#v\n", inputBuf)
}
```

#### Observe the specified key and decode

```go
func main() {
	// define callback function to process the data being observed
	callback := func(v []byte) (interface{}, error) {
		return y3.ToFloat32(v)
	}
	// create the Observable interface
	source := y3.FromStream(bytes.NewReader(inputBuf))
  // subscribe the Key being observed and set the callback function
	consumer := source.Subscribe(0x11).OnObserve(callback)
	// checking data after it has been processed
	for c := range consumer {
		fmt.Printf("observed value=%v, type=%v\n", c, reflect.ValueOf(c).Kind())
	}
}
```

## Types

Y3 provides support for decoding base types and structs and their slices, which can be done by the following methods：

<details>
  <summary>y3.ToObject</summary>
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
  <summary>y3.ToObject Slice</summary>
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
  <summary>y3.ToInt32</summary>
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
  <summary>y3.ToInt32Slice</summary>
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
  <summary>y3.ToUInt32</summary>
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
  <summary>y3.ToUInt32Slice</summary>
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
  <summary>y3.ToInt64</summary>
  <pre class="go" style="background-color: aliceblue">
	// Simulate source to generate and send data
	data := []<b>int64</b>{123, 456}
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
  <summary>y3.ToInt64Slice</summary>
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
  <summary>y3.ToUInt64</summary>
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
  <summary>y3.ToUInt64Slice</summary>
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
  <summary>y3.ToFloat32</summary>
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
  <summary>y3.ToFloat32Slice</summary>
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
  <summary>y3.ToFloat64</summary>
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
  <summary>y3.ToFloat64Slice</summary>
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
  <summary>y3.ToBool</summary>
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
  <summary>y3.ToBoolSlice</summary>
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
  <summary>y3.ToUTF8String</summary>
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
  <summary>y3.ToUTF8StringSlice</summary>
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

More examples in `/examples/`

## Attention

The key being subscribed has a defined range: 

- User-defined key to be subscribed: `0x10 ~ 0x3f`
- System reserved: `0x01 ~ 0x0f`

