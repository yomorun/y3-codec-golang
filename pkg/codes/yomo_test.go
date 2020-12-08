package codes

import (
	"fmt"
	"io"
	"testing"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"
)

func TestReadThermometer(t *testing.T) {
	inputBuf := buildThermometerInputData()

	codec := NewCodec("0x20")
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

func buildThermometerInputData() []byte {
	input := ThermometerTest{
		Id:          "the0",
		Temperature: float32(64.88),
		Humidity:    float32(93.02),
		Stored:      true,
	}

	proto := NewProtoCodec(packetutils.KeyOf("0x20"))
	inputBuf, _ := proto.Marshal(input)

	return inputBuf
}

func handleThermometer(value interface{}) (interface{}, error) {
	the := value.(*ThermometerTest)
	the.Id = "the1"
	return the, nil
}

func TestReadThermometerSlice(t *testing.T) {
	inputBuf := buildThermometerSliceInputData()

	codec := NewCodec("0x20")
	codec.Decoder(inputBuf)

	var mold = &[]ThermometerTest{}
	val, err := codec.Read(mold)
	if err != nil {
		fmt.Printf("#70 err1=%v\n", err)
		_ = val
	}

	result, err := handleThermometerSlice(mold)
	if err != nil {
		fmt.Printf("#70 err2=%v\n", err)
	}
	//fmt.Printf("#70 result=%v\n", result)

	res := result.([]ThermometerTest)
	if len(res) != 3 {
		t.Errorf("new result len should be: 3")
	}

	last := res[len(res)-1]
	if last.Id != "the1" {
		t.Errorf("last Id value should be: the1")
	}
	if last.Temperature != 1 {
		t.Errorf("last Temperature value should be: 1")
	}
	if last.Humidity != 2 {
		t.Errorf("last Humidity value should be: 2")
	}
	if last.Stored != false {
		t.Errorf("last Stored value should be: false")
	}
}

func handleThermometerSlice(value interface{}) (interface{}, error) {
	the := *value.(*[]ThermometerTest)
	the = append(the, ThermometerTest{
		Id:          "the1",
		Temperature: float32(1),
		Humidity:    float32(2),
		Stored:      false,
	})
	return the, nil
}

func buildThermometerSliceInputData() []byte {
	input := []ThermometerTest{
		{
			Id:          "the0",
			Temperature: float32(64.88),
			Humidity:    float32(93.02),
			Stored:      true,
		},
		{
			Id:          "the1",
			Temperature: float32(50),
			Humidity:    float32(90),
			Stored:      true,
		},
	}

	proto := NewProtoCodec(packetutils.KeyOf("0x20"))
	inputBuf, _ := proto.Marshal(input)

	return inputBuf
}

func TestComplexData(t *testing.T) {
	inputBuf := utils.DefaultTestedComplexData()

	// Decoder
	codec := NewCodec("0x12")
	codec.Decoder(inputBuf)

	// Read
	var mold = &utils.TestedResponseData{}
	val, err := codec.Read(mold)
	if err != nil {
		t.Errorf("Read error: %v", err)
	}

	if val == nil {
		t.Errorf("Read val is nil")
	}

	responseData := val.(*utils.TestedResponseData)
	if responseData.DateTime != "2017-09-11 10:52:27.811321" {
		t.Errorf("DateTime is unequal. DateTime=%v\n", responseData.DateTime)
	}

	searchData := responseData.SearchData
	if len(searchData) != 4 {
		t.Errorf("SearchData len is unequal. len(searchData)=%v\n", len(searchData))
	}

	lastContent := searchData[len(searchData)-1].Content
	lastUser := lastContent[len(lastContent)-1].User
	if lastUser.Id != 2384288641 {
		t.Errorf("lastUser Id is unequal. Id=%v\n", lastUser.Id)
	}

	experience := lastUser.Experience
	if experience.LevelInfo.Value != 1 {
		t.Errorf("experience LevelInfo Value is unequal. Value=%v\n", experience.LevelInfo.Value)
	}

	// process
	result, _ := process(responseData)

	// Write
	_, err = codec.Write(&complexDataWriter{}, result, mold)
	if err != nil {
		t.Errorf("Write error:%v", err)
	}
}

func process(value interface{}) (v interface{}, e error) {
	data := value.(*utils.TestedResponseData)
	content := data.SearchData[len(data.SearchData)-1].Content
	content[len(content)-1].User.Experience.LevelInfo.Name = "info"
	return *data, nil
}

type complexDataWriter struct{ io.Writer }

func (w *complexDataWriter) Write(buf []byte) (int, error) {
	// Decoder
	codec := NewCodec("0x12")
	codec.Decoder(buf)
	// Read
	var mold = &utils.TestedResponseData{}
	val, err := codec.Read(mold)
	if err != nil {
		panic(fmt.Errorf("read error: %v", err))
	}
	if val == nil {
		panic(fmt.Errorf("read val is nil"))
	}

	data := val.(*utils.TestedResponseData)
	content := data.SearchData[len(data.SearchData)-1].Content
	name := content[len(content)-1].User.Experience.LevelInfo.Name
	if name != "info" {
		panic(fmt.Errorf("User.Experience.LevelInfo.Name is unequal. name=%v\n", name))
	}

	return 0, nil
}
