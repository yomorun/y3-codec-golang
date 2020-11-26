package codes

import (
	"testing"
)

type ThermometerTest struct {
	Id          string  `yomo:"0x10"`
	Temperature float32 `yomo:"0x11"`
	Humidity    float32 `yomo:"0x12"`
	Stored      bool    `yomo:"0x13"`
}

func TestUnmarshalStructThermometerSlice(t *testing.T) {
	input := []ThermometerTest{
		{
			Id:          "the0",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
		{
			Id:          "the1",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
	}

	proto1 := NewProtoCodec("0x20")
	inputBuf, _ := proto1.Marshal(input)

	proto2 := NewProtoCodec("0x20")
	var mold []ThermometerTest
	runUnmarshalStruct(t, proto2, inputBuf, &mold)

	if len(mold) != len(input) {
		t.Errorf("len value should be: %v", len(input))
	}

	for i, v := range mold {
		testThermometerStruct(t, v, input[i])
	}
}

func TestUnmarshalStructThermometer(t *testing.T) {
	input := ThermometerTest{
		Id:          "the0",
		Temperature: float32(64.88),
		Humidity:    float32(93.02),
		Stored:      true,
	}

	proto1 := NewProtoCodec("0x20")
	inputBuf, _ := proto1.Marshal(input)

	//debug:
	//fmt.Printf("#60 buf=%s\n", packetutils.FormatBytes(inputBuf))
	//fmt.Printf("#60 node=")
	//res, _, _ := y3.DecodeNodePacket(inputBuf)
	//packetutils.PrintNodePacket(res)
	//fmt.Println()

	proto2 := NewProtoCodec("0x20")
	var mold = ThermometerTest{}
	runUnmarshalStruct(t, proto2, inputBuf, &mold)
	testThermometerStruct(t, mold, input)
}

func runUnmarshalStruct(t *testing.T, proto ProtoCodec, inputBuf []byte, mold interface{}) {
	err := proto.UnmarshalStruct(inputBuf, mold)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}
}

func testThermometerStruct(t *testing.T, result ThermometerTest, expected ThermometerTest) {
	if result.Id != expected.Id {
		t.Errorf("Id value should be: %v", expected.Id)
	}

	if result.Temperature != expected.Temperature {
		t.Errorf("Temperature value should be: %v", expected.Temperature)
	}

	if result.Humidity != expected.Humidity {
		t.Errorf("Humidity value should be: %v", expected.Humidity)
	}

	if result.Stored != expected.Stored {
		t.Errorf("Stored value should be: %v", expected.Stored)
	}
}
