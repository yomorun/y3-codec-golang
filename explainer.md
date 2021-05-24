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

Y3 provides `High-Level` wrappers for supporting frameworks like [YoMo](https://github.com/yomorun/yomo).

Unified encoding method: `y3.NewCodec(observe byte).Marshal(input interface{})`

| Types         | Decoding Methods     |
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
  
  ```golang
  func main() {
    // Simulate source to generate and send data
    data := NoiseData{Noise: 40, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"}
    sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
    source := y3.FromStream(bytes.NewReader(sendingBuf))
    // Simulate flow listening and decoding data
    var decode = func(v []byte) (interface{}, error) {
        var obj NoiseData
        err := y3.ToObject(v, &obj)
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
  type NoiseData struct {
      Noise float32 `y3:"0x11"`
      Time  int64   `y3:"0x12"`
      From  string  `y3:"0x13"`
  }
  ```
</details>
<details>
  <summary>struct slice</summary>
  
  ```golang
    func main() {
      // Simulate source to generate and send data
      data := []NoiseData{
          {Noise: 40, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
          {Noise: 50, Time: time.Now().UnixNano() / 1e6, From: "127.0.0.1"},
      }
      sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
      source := y3.FromStream(bytes.NewReader(sendingBuf))
      // Simulate flow listening and decoding data
      var decode = func(v []byte) (interface{}, error) {
          var sl []NoiseData
          err := y3.ToObject(v, &sl)
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
    type NoiseData struct {
        Noise float32 `y3:"0x11"`
        Time  int64   `y3:"0x12"`
        From  string  `y3:"0x13"`
    }
  ```
</details>
<details>
  <summary>int32</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data int32 = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToInt32(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>int32 slice</summary>
  
  ```golang
    // Simulate source to generate and send data
    data := []int32{123, 456}
    sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
    source := y3.FromStream(bytes.NewReader(sendingBuf))
    // Simulate flow listening and decoding data
    var decode = func(v []byte) (interface{}, error) {
        sl, err := y3.ToInt32Slice(v)
        if err != nil {
            return nil, err
        }
        fmt.Printf("encoded data: %v\n", sl)
        return sl, nil
    }
    consumer := source.Subscribe(0x10).OnObserve(decode)
    for range consumer {
    }
  ```
</details>
<details>
  <summary>uint32</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data uint32 = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToUInt32(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>uint32 slice</summary>
  
  ```golang
  // Simulate source to generate and send data
  data := []uint32{123, 456}
  sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
  source := y3.FromStream(bytes.NewReader(sendingBuf))
  // Simulate flow listening and decoding data
  var decode = func(v []byte) (interface{}, error) {
      sl, err := y3.ToUInt32Slice(v)
      if err != nil {
          return nil, err
      }
      fmt.Printf("encoded data: %v\n", sl)
      return sl, nil
  }
  consumer := source.Subscribe(0x10).OnObserve(decode)
  for range consumer {
  }
  ```
</details>
<details>
  <summary>int64</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data int64 = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToInt64(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>int64 slice</summary>
  
  ```golang
  // Simulate source to generate and send data
  data := []int64{123, 456}
  sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
  source := y3.FromStream(bytes.NewReader(sendingBuf))
  // Simulate flow listening and decoding data
  var decode = func(v []byte) (interface{}, error) {
      sl, err := y3.ToInt64Slice(v)
      if err != nil {
          return nil, err
      }
      fmt.Printf("encoded data: %v\n", sl)
      return sl, nil
  }
  consumer := source.Subscribe(0x10).OnObserve(decode)
  for range consumer {
  }
  ```
</details>
<details>
  <summary>uint64</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data uint64 = 123
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToUInt64(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>uint64 slice</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := []uint64{123, 456}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToUInt64Slice(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>float32</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data float32 = 1.23
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToFloat32(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>float32 slice</summary>
  
  ```golang
  // Simulate source to generate and send data
	data := []float32{1.23, 4.56}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToFloat32Slice(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>float64</summary>
  
  ```golang
	// Simulate source to generate and send data
	var data float64 = 1.23
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToFloat64(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>float64 slice</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := []float64{1.23, 4.56}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToFloat64Slice(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>bool</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := true
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToBool(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>bool slice</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := []bool{true, false}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToBoolSlice(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>string</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := "abc"
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToUTF8String(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>
<details>
  <summary>string slice</summary>
  
  ```golang
	// Simulate source to generate and send data
	data := []string{"a", "b"}
	sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)
	source := y3.FromStream(bytes.NewReader(sendingBuf))
	// Simulate flow listening and decoding data
	var decode = func(v []byte) (interface{}, error) {
		sl, err := y3.ToUTF8StringSlice(v)
		if err != nil {
			return nil, err
		}
		fmt.Printf("encoded data: %v\n", sl)
		return sl, nil
	}
	consumer := source.Subscribe(0x10).OnObserve(decode)
	for range consumer {
	}
  ```
</details>

More examples in `/examples/`

## Attention

The key being subscribed has a defined range: 

- User-defined key to be subscribed: `0x10 ~ 0x3f`
- System reserved: `0x01 ~ 0x0f`

