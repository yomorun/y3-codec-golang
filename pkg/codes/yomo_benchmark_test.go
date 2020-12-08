package codes

import (
	"fmt"
	"testing"
)

func BenchmarkBasicMax(b *testing.B) {
	inputBuf := []byte{0x80, 0x81, 0x10, 0x29, 0x1, 0x29, 0x30, 0x1, 0x30, 0x3f, 0x1, 0x3f, 0x1e, 0x1, 0x1e, 0x1d, 0x1, 0x1d, 0x23, 0x1, 0x23, 0x16, 0x1, 0x16, 0x18, 0x1, 0x18, 0x1f, 0x1, 0x1f, 0x2a, 0x1, 0x2a, 0x2c, 0x1, 0x2c, 0x33, 0x1, 0x33, 0x37, 0x1, 0x37, 0x15, 0x1, 0x15, 0x11, 0x1, 0x11, 0x12, 0x1, 0x12, 0x13, 0x1, 0x13, 0x1c, 0x1, 0x1c, 0x2b, 0x1, 0x2b, 0x2f, 0x1, 0x2f, 0x3b, 0x1, 0x3b, 0x10, 0x1, 0x10, 0x26, 0x1, 0x26, 0x28, 0x1, 0x28, 0x32, 0x1, 0x32, 0x36, 0x1, 0x36, 0x38, 0x1, 0x38, 0x3c, 0x1, 0x3c, 0x1a, 0x1, 0x1a, 0x21, 0x1, 0x21, 0x22, 0x1, 0x22, 0x24, 0x1, 0x24, 0x25, 0x1, 0x25, 0x2d, 0x1, 0x2d, 0x2e, 0x1, 0x2e, 0x31, 0x1, 0x31, 0x17, 0x1, 0x17, 0x3d, 0x1, 0x3d, 0x3e, 0x1, 0x3e, 0x3a, 0x1, 0x3a, 0x20, 0x1, 0x20, 0x27, 0x1, 0x27, 0x34, 0x1, 0x34, 0x39, 0x1, 0x39, 0x1b, 0x1, 0x1b, 0x19, 0x1, 0x19, 0x35, 0x1, 0x35, 0x14, 0x1, 0x14}
	codec := NewCodec("0x23")
	//codec.Decoder(inputBuf)
	//v, e := codec.Read(uint32(1))
	//fmt.Printf("v=%v, e=%v\n", v, e)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Decoder(inputBuf)
		if _, err := codec.Read(uint32(1)); err != nil {
			panic(fmt.Errorf("read is failure: %v", err))
		}
	}
}

func BenchmarkBasic(b *testing.B) {
	inputBuf := []byte{0x80, 0xe, 0x11, 0x1, 0x19, 0x10, 0x9, 0x7a, 0x68, 0x61, 0x6e, 0x67, 0x20, 0x73, 0x61, 0x6e}
	codec := NewCodec("0x10")
	//codec.Decoder(inputBuf)
	//v, e := codec.Read("")
	//fmt.Printf("v=%v, e=%v\n", v, e)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Decoder(inputBuf)
		if _, err := codec.Read(""); err != nil {
			panic(fmt.Errorf("read is failure: %v", err))
		}
	}
}

func BenchmarkPerson(b *testing.B) {
	observe := "0x20"
	inputBuf := NewCodecBenchmarkData().DefaultPersonData()
	codec := NewCodec(observe)
	//codec.Decoder(inputBuf)
	//v, e := codec.Read(&Person{})
	//fmt.Printf("v=%v, e=%v\n", v, e)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Decoder(inputBuf)
		if _, err := codec.Read(&Person{}); err != nil {
			panic(fmt.Errorf("read is failure: %v", err))
		}
	}
}

func BenchmarkPersonMax(b *testing.B) {
	inputBuf := NewCodecBenchmarkData().DefaultPersonMaxData()
	codec := NewCodec("0x23")
	//codec.Decoder(inputBuf)
	//v, e := codec.Read(&Person{})
	//fmt.Printf("v=%v, e=%v\n", v, e)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Decoder(inputBuf)
		if _, err := codec.Read(&Person{}); err != nil {
			panic(fmt.Errorf("read is failure: %v", err))
		}
	}
}
