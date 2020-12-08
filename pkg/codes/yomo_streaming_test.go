package codes

import (
	"testing"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"
)

// TL,V
func TestCodecDecoderStick(t *testing.T) {
	//inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	inputBuf1 := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9}
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	inputBuf3 := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9}
	inputBuf4 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	//fmt.Printf("inputBuf=%v", packetutils.FormatBytes(inputBuf1))
	//fmt.Printf(" %v\n", packetutils.FormatBytes(inputBuf2))

	codec := NewStreamingCodec(0x10) //0x11,0x10
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)
	codec.Decoder(inputBuf3)
	codec.Decoder(inputBuf4)
}

func TestCodecDecoder(t *testing.T) {
	inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	codec := NewStreamingCodec(0x10) //0x11,0x10
	codec.Decoder(inputBuf)
}

func TestCodecReadBasic(t *testing.T) {
	inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	codec := NewStreamingCodec(0x10) //0x11,0x10
	codec.Decoder(inputBuf)

	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	_ = val
	//fmt.Printf("val=%v\n", val)

}

func TestCodecReadThermometer(t *testing.T) {
	inputBuf := buildThermometerInputData()
	codec := NewStreamingCodec(0x20)
	codec.Decoder(inputBuf)

	var mold = &ThermometerTest{}
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=%v\n", val)

	if val == nil {
		t.Errorf("Read val is nil")
	}

	result, err := handleThermometer(val)
	if err != nil {
		t.Errorf("handleThermometer error: %v", err)
	}

	if result.(*ThermometerTest).Id != "the1" {
		t.Errorf("Id value should be: the1")
	}
}

func TestCodecReadThermometerMax(t *testing.T) {
	inputBuf := NewCodecBenchmarkData().DefaultPersonMaxData()
	codec := NewStreamingCodec(packetutils.KeyOf("0x23"))
	codec.Decoder(inputBuf)

	var mold = &ThermometerTest{}
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}

	if val == nil {
		t.Errorf("Read val is nil")
	}

	result, err := handleThermometer(val)
	if err != nil {
		t.Errorf("handleThermometer error: %v", err)
	}

	if result.(*ThermometerTest).Id != "the1" {
		t.Errorf("Id value should be: the1")
	}
}
