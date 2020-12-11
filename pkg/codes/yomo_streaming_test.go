package codes

import (
	"math/rand"
	"testing"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"
)

func TestCodecDecoderStickyTLV_RandomCut_Name(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBufTwoCopy := append(inputBuf, inputBuf...)
	point := rand.Int31n(int32(len(inputBufTwoCopy)))
	inputBuf1 := inputBufTwoCopy[:point]
	inputBuf2 := inputBufTwoCopy[point:]

	codec := NewStreamingCodec(0x10) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	// read1
	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}

	// read2
	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}
}

func TestCodecDecoderStickyTLV_RandomCut_Birthday(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBufTwoCopy := append(inputBuf, inputBuf...)
	point := rand.Int31n(int32(len(inputBufTwoCopy)))
	inputBuf1 := inputBufTwoCopy[:point]
	inputBuf2 := inputBufTwoCopy[point:]

	codec := NewStreamingCodec(0x12) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	// read1
	var mold int64
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}

	// read2
	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}
}

func TestCodecDecoderStickyTLV_OneByOne_Name(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	realBuf := append(inputBuf, inputBuf...)

	codec := NewStreamingCodec(0x10) //0x10,0x11,0x12
	for _, c := range realBuf {
		codec.Decoder([]byte{c})
	}

	// read1
	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}

	// read2
	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}
}

func TestCodecDecoderStickyTLV_OneByOne_Birthday(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	realBuf := append(inputBuf, inputBuf...)

	codec := NewStreamingCodec(0x12) //0x10,0x11,0x12
	for _, c := range realBuf {
		codec.Decoder([]byte{c})
	}

	// read1
	var mold int64
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}

	// read2
	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}
}

func TestCodecDecoderStickyTLV_CutInSecondPacket_Name(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf1 := append([]byte{}, inputBuf...)
	inputBuf1 = append(inputBuf1, []byte{0x81, 0x16, 0x10, 0x9}...)
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}

	codec := NewStreamingCodec(0x10) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}
}

func TestCodecDecoderStickyTLV_CutInSecondPacket_Birthday(t *testing.T) {
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf1 := append([]byte{}, inputBuf...)
	inputBuf1 = append(inputBuf1, []byte{0x81, 0x16, 0x10, 0x9}...)
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}

	codec := NewStreamingCodec(0x12) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	var mold int64
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}
}

func TestCodecDecoderStickyTLV_OnTwoPacket_Name(t *testing.T) {
	//inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf1 := []byte{0x81, 0x16, 0x10, 0x9}
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}

	codec := NewStreamingCodec(0x10) //0x11,0x10
	//codec.Decoder(inputBuf)
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("Name value should be: zhang san")
	}
}

func TestCodecDecoderStickyTLV_OnTwoPacket_Name2(t *testing.T) {
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

	// read1
	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("0x10 value should be: zhang san")
	}

	// read2
	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != "zhang san" {
		t.Errorf("0x10 value should be: zhang san")
	}
}

func TestCodecDecoderStickyTLV_OnTwoPacket_Age(t *testing.T) {
	inputBuf1 := []byte{0x81, 0x16, 0x10, 0x9}
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}

	codec := NewStreamingCodec(0x11) //0x11,0x10
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	var mold int32
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int32(255) {
		t.Errorf("Age value should be: 255")
	}
}

func TestCodecDecoderStickyTLV_OnTwoPacket_Birthday(t *testing.T) {
	//inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf1 := []byte{0x81, 0x16, 0x10, 0x9}
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}

	codec := NewStreamingCodec(0x12) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)
	//codec.Decoder(inputBuf)

	var mold int64
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}
}

func TestCodecDecoderStickyTLV_CutInFirstPacket_Birthday(t *testing.T) {
	//inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf := []byte{0x81, 0x16, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf1 := []byte{0x81, 0x16, 0x10, 0x9}
	inputBuf2 := []byte{0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e, 0x11, 0x2, 0x81, 0x7f, 0x12, 0x5, 0x83, 0xc3, 0xb3, 0xa6, 0x0}
	inputBuf2 = append(inputBuf2, inputBuf...)

	codec := NewStreamingCodec(0x12) //0x10,0x11,0x12
	codec.Decoder(inputBuf1)
	codec.Decoder(inputBuf2)

	var mold int64
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}

	val, err = codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=`%v`\n", val)
	if val != int64(946656000) {
		t.Errorf("Birthday value should be: 946656000")
	}
}

func TestCodecReadBasic_Age(t *testing.T) {
	inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	codec := NewStreamingCodec(0x11) //0x11,0x10
	codec.Decoder(inputBuf)

	var mold int32
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=%v\n", val)

	if val != int32(25) {
		t.Errorf("0x11 value should be: 25")
	}
}

func TestCodecReadBasic_Name(t *testing.T) {
	inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	codec := NewStreamingCodec(0x10) //0x11,0x10
	codec.Decoder(inputBuf)

	var mold = ""
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}
	//fmt.Printf("val=%v\n", val)

	if val != "zhang san" {
		t.Errorf("0x10 value should be: zhang san")
	}

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

	if val.(*ThermometerTest).Humidity != 93.02 {
		t.Errorf("Humidity value should be: 93.02")
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
