package main

import (
	"fmt"
	"time"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	"github.com/yomorun/yomo-codec-golang/internal/pprof"
	"github.com/yomorun/yomo-codec-golang/pkg/codes"
)

func main() {
	observe := "0x20"
	inputBuf := buildPersonData(observe)

	// pprof
	fmt.Printf("start pprof\n")
	go pprof.Run()
	time.Sleep(5 * time.Second)

	fmt.Printf("start testing...\n")
	//codec := codes.NewCodec(observe)
	//for {
	//	codec.Decoder(inputBuf)
	//	if _, err := codec.Read(&Person{}); err != nil {
	//		panic(fmt.Errorf("read is failure: %v", err))
	//	}
	//	//time.Sleep(200 * time.Microsecond)
	//}
	// #V2
	codec := codes.NewStreamingCodec(packetutils.KeyOf(observe))
	for {
		codec.Decoder(inputBuf)
		if _, err := codec.Read(&Person{}); err != nil {
			panic(fmt.Errorf("read is failure: %v", err))
		}
		//time.Sleep(200 * time.Microsecond)
	}

	//time.Sleep(3600 * time.Second)
}

type Person struct {
	Name string `yomo:"0x10"`
	Age  uint32 `yomo:"0x11"`
}

func buildPersonData(observe string) []byte {
	input := Person{
		Name: "zhang san",
		Age:  25,
	}

	proto := codes.NewProtoCodec(packetutils.KeyOf(observe))
	inputBuf, _ := proto.Marshal(input)

	return inputBuf
}
