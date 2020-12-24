package codes

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"
)

var (
	testLoopTimes = 10
)

type InformSample struct {
	Name string `yomo:"0x11"`
}

func TestInformBySample(t *testing.T) {
	source := InformSample{Name: "yomo"}
	target := InformSample{Name: "yomo!"}

	input, _ := NewProtoCodec(0x10).Marshal(source)

	handle := func(v interface{}) interface{} {
		the := v.(*InformSample)
		the.Name = "yomo!"
		return *the
	}

	verify := func(buf []byte) {
		var mold InformSample
		NewProtoCodec(0x10).UnmarshalStruct(buf, &mold)
		if mold.Name != target.Name {
			t.Errorf("TestInformBySample, The value should be %v", target.Name)
		}
	}

	testInformBy(t, &InformSample{}, input, handle, verify)
}

func TestInformBySampleSlice(t *testing.T) {
	source := []InformSample{{Name: "yomo"}}
	target := []InformSample{{Name: "yomo!"}}

	input, _ := NewProtoCodec(0x10).Marshal(source)

	handle := func(v interface{}) interface{} {
		the := *v.(*[]InformSample)
		for i, _ := range the {
			the[i].Name = "yomo!"
		}
		return the
	}

	verify := func(buf []byte) {
		var mold []InformSample
		NewProtoCodec(0x10).UnmarshalStruct(buf, &mold)
		if mold[0].Name != target[0].Name {
			t.Errorf("TestInformBySampleSlice, The value should be %v", target[0].Name)
		}
	}

	testInformBy(t, &[]InformSample{}, input, handle, verify)
}

func TestInformByString(t *testing.T) {
	source := "y-new"
	target := "y-new!"

	input, _ := NewProtoCodec(0x10).Marshal(source)

	handle := func(v interface{}) interface{} {
		return v.(string) + "!"
	}

	verify := func(buf []byte) {
		var mold interface{} = ""
		NewProtoCodec(0x10).UnmarshalBasic(buf, &mold)
		if mold != "y-new!" {
			t.Errorf("TestInformByString, The value should be %v", target)
		}
	}

	testInformBy(t, "", input, handle, verify)
}

func TestInformByStringSlice(t *testing.T) {
	source := []string{"a", "b", "c"}
	target := []string{"a", "b", "c", "d"}

	input, _ := NewProtoCodec(0x10).Marshal(source)

	handle := func(v interface{}) interface{} {
		if v == nil {
			panic(fmt.Errorf("#444 handleStringSlice v cannot be nil"))
		}

		value := v.([]interface{})
		result := make([]string, 0)
		for _, v := range value {
			result = append(result, fmt.Sprintf("%v", v))
		}
		result = append(result, "d")
		return result
	}

	verify := func(buf []byte) {
		var mold interface{} = []string{}
		NewProtoCodec(0x10).UnmarshalBasic(buf, &mold)
		arr, _ := utils.ToStringSliceArray(mold)
		for i, v := range arr {
			if v != target[i] {
				t.Errorf("TestInformByStringSlice, The value should be %v=%v", i, target[i])
			}
		}

	}

	testInformBy(t, []string{}, input, handle, verify)
}

func testInformBy(t *testing.T, mold interface{}, input []byte, handle func(v interface{}) interface{}, verify func(buf []byte)) {
	codec, inform := NewStreamingCodec(0x10)
	go func() {
		l := testLoopTimes
		for i := 0; i < l; i++ {
			codec.Decoder(input)
		}
	}()

	for {
		select {
		case flag := <-inform:
			if flag {
				v, _ := codec.Read(mold)
				vv := handle(v)
				codec.Write(&verifier{Verify: verify}, vv, mold)
			}
			codec.Refresh(&verifier{Verify: verify})
		case <-time.After(1 * time.Second):
			fmt.Printf("%v select timeout to quit\n", time.Now().Format("2006-01-02 15:04:05"))
			return
		}
	}
}

type verifier struct {
	Verify func(buf []byte)
	io.Writer
}

func (w *verifier) Write(buf []byte) (n int, err error) {
	fmt.Printf("%v %v\n", time.Now().Format("2006-01-02 15:04:05"), packetutils.FormatBytes(buf))
	w.Verify(buf)
	return len(buf), nil
}
