package y3

import (
	"reflect"
	"testing"

	"github.com/yomorun/y3-codec-golang/internal/utils"

	"github.com/yomorun/y3-codec-golang/internal/tester"
)

func newBasic() tester.BasicTestData {
	return tester.BasicTestData{
		Vstring:  "foo",
		Vint32:   int32(127),
		Vint64:   int64(-1),
		Vuint32:  uint32(130),
		Vuint64:  uint64(18446744073709551615),
		Vfloat32: float32(0.25),
		Vfloat64: float64(23),
		Vbool:    true,
	}
}

func TestBasic_Struct(t *testing.T) {
	t.Parallel()

	input := newBasic()
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.BasicTestData
	runDecode(t, inputBuf, &result)
	testBasicStruct(t, result, input)
}

func runDecode(t *testing.T, inputBuf []byte, output interface{}) {
	_, err := newStructDecoder(output, structDecoderOptionConfig(&structDecoderConfig{
		ZeroFields: true,
		TagName:    "y3",
	})).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}
}

func TestDecode_Embedded(t *testing.T) {
	t.Parallel()

	input := tester.EmbeddedTestData{
		BasicTestData: newBasic(),
		Vaction:       "drink",
	}
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.EmbeddedTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.Vaction != "drink" {
		t.Errorf("vstring value should be 'drink': %#v", result.Vaction)
	}

	testBasicStruct(t, result.BasicTestData, input.BasicTestData)
}

func TestDecode_EmbeddedMore(t *testing.T) {
	t.Parallel()

	input := tester.EmbeddedMoreTestData{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"}
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.EmbeddedMoreTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.Vanimal != "bird" {
		t.Errorf("vstring value should be 'bird': %#v", result.Vanimal)
	}

	if result.Vaction != "drink" {
		t.Errorf("vstring value should be 'drink': %#v", result.Vaction)
	}

	testBasicStruct(t, result.BasicTestData, input.BasicTestData)
}

func TestDecoder_Named(t *testing.T) {
	t.Parallel()

	input := tester.NamedTestData{Base: newBasic(), Vaction: "drink"}
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.NamedTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.Vaction != "drink" {
		t.Errorf("vstring value should be 'drink': %#v", result.Vaction)
	}

	testBasicStruct(t, result.Base, input.Base)
}

func TestDecoder_NamedMore(t *testing.T) {
	t.Parallel()

	input := tester.NamedMoreTestData{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.NamedMoreTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.Vanimal != "bird" {
		t.Errorf("vstring value should be 'bird': %#v", result.Vanimal)
	}

	if result.MyNest.Vaction != "drink" {
		t.Errorf("vstring value should be 'drink': %#v", result.MyNest.Vaction)
	}

	testBasicStruct(t, result.MyNest.Base, input.MyNest.Base)
}

func TestArray(t *testing.T) {
	t.Parallel()

	input := tester.ArrayTestData{
		Vfoo:          "foo",
		Vbar:          [2]string{"foo", "bar"},
		Vint32Array:   [2]int32{1, 2},
		Vint64Array:   [2]int64{1, 2},
		Vuint32Array:  [2]uint32{1, 2},
		Vuint64Array:  [2]uint64{1, 2},
		Vfloat32Array: [2]float32{1, 2},
		Vfloat64Array: [2]float64{1, 2},
	}
	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.ArrayTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testArrayString(t, result.Vbar, input.Vbar)
	testArrayInt32(t, result.Vint32Array, input.Vint32Array)
	testArrayInt64(t, result.Vint64Array, input.Vint64Array)
	testArrayUint32(t, result.Vuint32Array, input.Vuint32Array)
	testArrayUint64(t, result.Vuint64Array, input.Vuint64Array)
	testArrayFloat32(t, result.Vfloat32Array, input.Vfloat32Array)
	testArrayFloat64(t, result.Vfloat64Array, input.Vfloat64Array)
}

func TestSlice(t *testing.T) {
	t.Parallel()

	input := tester.SliceTestData{
		Vfoo:          "foo",
		Vbar:          []string{"foo", "bar"},
		Vint32Slice:   []int32{1, 2},
		Vint64Slice:   []int64{1, 2},
		Vuint32Slice:  []uint32{1, 2},
		Vuint64Slice:  []uint64{1, 2},
		Vfloat32Slice: []float32{1, 2},
		Vfloat64Slice: []float64{1, 2},
	}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.SliceTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testSliceString(t, result.Vbar, input.Vbar)
	testSliceInt32(t, result.Vint32Slice, input.Vint32Slice)
	testSliceInt64(t, result.Vint64Slice, input.Vint64Slice)
	testSliceUint32(t, result.Vuint32Slice, input.Vuint32Slice)
	testSliceUint64(t, result.Vuint64Slice, input.Vuint64Slice)
	testSliceFloat32(t, result.Vfloat32Slice, input.Vfloat32Slice)
	testSliceFloat64(t, result.Vfloat64Slice, input.Vfloat64Slice)
}

func TestEmptySlice(t *testing.T) {
	t.Parallel()

	input := tester.SliceTestData{
		Vfoo:          "foo",
		Vbar:          []string{},
		Vint32Slice:   []int32{},
		Vint64Slice:   []int64{},
		Vuint32Slice:  []uint32{},
		Vuint64Slice:  []uint64{},
		Vfloat32Slice: []float32{},
		Vfloat64Slice: []float64{},
	}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.SliceTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testSliceString(t, result.Vbar, input.Vbar)
	testSliceInt32(t, result.Vint32Slice, input.Vint32Slice)
	testSliceInt64(t, result.Vint64Slice, input.Vint64Slice)
	testSliceUint32(t, result.Vuint32Slice, input.Vuint32Slice)
	testSliceUint64(t, result.Vuint64Slice, input.Vuint64Slice)
	testSliceFloat32(t, result.Vfloat32Slice, input.Vfloat32Slice)
	testSliceFloat64(t, result.Vfloat64Slice, input.Vfloat64Slice)
}

func TestSliceStruct(t *testing.T) {
	t.Parallel()

	input := tester.SliceStructTestData{
		Vstring:  "foo",
		BaseList: []tester.BasicTestData{newBasic(), newBasic()},
		NamedMoreList: []tester.NamedMoreTestData{
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: []tester.EmbeddedMoreTestData{
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.SliceStructTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValue(t, result.Vstring, input.Vstring)

	testValueWith(t, len(result.BaseList), len(input.BaseList), "BaseList len")

	for i, v := range result.BaseList {
		testBasicStruct(t, v, input.BaseList[i])
	}

	testValueWith(t, len(result.NamedMoreList), len(input.NamedMoreList), "NamedMoreList len")

	for i, v := range result.NamedMoreList {
		testValue(t, v.Vanimal, input.NamedMoreList[i].Vanimal)
		testValue(t, v.MyNest.Vaction, input.NamedMoreList[i].MyNest.Vaction)
		testBasicStruct(t, v.MyNest.Base, input.NamedMoreList[i].MyNest.Base)
	}

	testValueWith(t, len(result.EmbeddedMoreList), len(input.EmbeddedMoreList), "EmbeddedMoreList len")

	for i, v := range result.EmbeddedMoreList {
		testValue(t, v.Vanimal, input.EmbeddedMoreList[i].Vanimal)
		testValue(t, v.Vaction, input.EmbeddedMoreList[i].Vaction)
		testBasicStruct(t, v.BasicTestData, input.EmbeddedMoreList[i].BasicTestData)
	}
}

func TestArrayStruct(t *testing.T) {
	t.Parallel()

	input := tester.ArrayStructTestData{
		Vstring:  "foo",
		BaseList: [2]tester.BasicTestData{newBasic(), newBasic()},
		NamedMoreList: [2]tester.NamedMoreTestData{
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: [2]tester.EmbeddedMoreTestData{
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.ArrayStructTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValue(t, result.Vstring, input.Vstring)

	testValueWith(t, len(result.BaseList), len(input.BaseList), "BaseList len")

	for i, v := range result.BaseList {
		testBasicStruct(t, v, input.BaseList[i])
	}

	testValueWith(t, len(result.NamedMoreList), len(input.NamedMoreList), "NamedMoreList len")

	for i, v := range result.NamedMoreList {
		testValue(t, v.Vanimal, input.NamedMoreList[i].Vanimal)
		testValue(t, v.MyNest.Vaction, input.NamedMoreList[i].MyNest.Vaction)
		testBasicStruct(t, v.MyNest.Base, input.NamedMoreList[i].MyNest.Base)
	}

	testValueWith(t, len(result.EmbeddedMoreList), len(input.EmbeddedMoreList), "EmbeddedMoreList len")

	for i, v := range result.EmbeddedMoreList {
		testValue(t, v.Vanimal, input.EmbeddedMoreList[i].Vanimal)
		testValue(t, v.Vaction, input.EmbeddedMoreList[i].Vaction)
		testBasicStruct(t, v.BasicTestData, input.EmbeddedMoreList[i].BasicTestData)
	}
}

func TestRootSliceWithBasicStruct(t *testing.T) {
	t.Parallel()

	input := []tester.BasicTestData{newBasic(), newBasic()}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result []tester.BasicTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValueWith(t, len(result), len(input), "[]BasicTestData len")

	for i, v := range result {
		testBasicStruct(t, v, input[i])
	}
}

func TestRootSliceWithSliceStruct(t *testing.T) {
	t.Parallel()

	input1 := tester.SliceStructTestData{
		Vstring:  "foo",
		BaseList: []tester.BasicTestData{newBasic(), newBasic()},
		NamedMoreList: []tester.NamedMoreTestData{
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: tester.NamedTestData{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: []tester.EmbeddedMoreTestData{
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{EmbeddedTestData: tester.EmbeddedTestData{BasicTestData: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	input := []tester.SliceStructTestData{input1, input1}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result []tester.SliceStructTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValueWith(t, len(result), len(input), "[]SliceStructTestData len")

	for i, v := range result {
		testValue(t, v.Vstring, input[i].Vstring)

		testValueWith(t, len(v.BaseList), len(input[i].BaseList), "BaseList len")

		for i, v := range v.BaseList {
			testBasicStruct(t, v, input[i].BaseList[i])
		}

		testValueWith(t, len(v.NamedMoreList), len(input[i].NamedMoreList), "NamedMoreList len")

		for i, vv := range v.NamedMoreList {
			testValue(t, vv.Vanimal, input[i].NamedMoreList[i].Vanimal)
			testValue(t, vv.MyNest.Vaction, input[i].NamedMoreList[i].MyNest.Vaction)
			testBasicStruct(t, vv.MyNest.Base, input[i].NamedMoreList[i].MyNest.Base)
		}

		testValueWith(t, len(v.EmbeddedMoreList), len(input[i].EmbeddedMoreList), "EmbeddedMoreList len")

		for i, vv := range v.EmbeddedMoreList {
			testValue(t, vv.Vanimal, input[i].EmbeddedMoreList[i].Vanimal)
			testValue(t, vv.Vaction, input[i].EmbeddedMoreList[i].Vaction)
			testBasicStruct(t, vv.BasicTestData, input[i].EmbeddedMoreList[i].BasicTestData)
		}
	}

}

func TestNested(t *testing.T) {
	t.Parallel()

	input := tester.NestedTestData{
		SubNested: tester.Sub1NestedTestData{
			SubNested: tester.Sub2NestedTestData{
				SubNested: tester.Sub3NestedTestData{
					BasicList: []tester.BasicTestData{newBasic(), newBasic()},
				}}}}

	inputBuf, _ := newStructEncoder(0x3f, structEncoderOptionRoot(utils.RootToken)).Encode(input)

	var result tester.NestedTestData
	_, err := newStructDecoder(&result).Decode(inputBuf)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	for i, v := range result.SubNested.SubNested.SubNested.BasicList {
		testBasicStruct(t, v, input.SubNested.SubNested.SubNested.BasicList[i])
	}

}

func testBasicStruct(t *testing.T, result tester.BasicTestData, expected tester.BasicTestData) {
	if result.Vstring != expected.Vstring {
		t.Errorf("Vstring value should be: %v", expected.Vstring)
	}

	if result.Vint32 != expected.Vint32 {
		t.Errorf("Vint32 value should be: %v", expected.Vint32)
	}

	if result.Vint64 != expected.Vint64 {
		t.Errorf("Vint64 value should be: %#v", expected.Vint64)
	}

	if result.Vuint32 != expected.Vuint32 {
		t.Errorf("Vuint32 value should be: %v", expected.Vuint32)
	}

	if result.Vuint64 != expected.Vuint64 {
		t.Errorf("Vuint64 value should be: %v", expected.Vuint64)
	}

	if result.Vfloat32 != expected.Vfloat32 {
		t.Errorf("Vfloat32 value should be: %v", expected.Vfloat32)
	}

	if result.Vfloat64 != expected.Vfloat64 {
		t.Errorf("Vfloat64 value should be: %v", expected.Vfloat64)
	}

	if result.Vbool != expected.Vbool {
		t.Errorf("Vbool value should be: %v", expected.Vbool)
	}
}

func testValue(t *testing.T, result interface{}, expected interface{}) {
	testValueWith(t, result, expected, "")
}

func testValueWith(t *testing.T, result interface{}, expected interface{}, prefix string) {
	getPrefix := func(prefix string, def string) string {
		if prefix != "" && len(prefix) > 0 {
			return prefix
		}
		return def
	}

	switch reflect.ValueOf(expected).Kind() {
	case reflect.String:
		if result.(string) != expected.(string) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "string"), expected)
		}
	case reflect.Int32:
		if result.(int32) != expected.(int32) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "int32"), expected)
		}
	case reflect.Int64:
		if result.(int64) != expected.(int64) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "int64"), expected)
		}
	case reflect.Uint32:
		if result.(uint32) != expected.(uint32) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "uint32"), expected)
		}
	case reflect.Uint64:
		if result.(uint64) != expected.(uint64) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "uint64"), expected)
		}
	case reflect.Float32:
		if result.(float32) != expected.(float32) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "float32"), expected)
		}
	case reflect.Float64:
		if result.(float64) != expected.(float64) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "float64"), expected)
		}
	case reflect.Bool:
		if result.(bool) != expected.(bool) {
			t.Errorf("%v value should be: %v", getPrefix(prefix, "bool"), expected)
		}
	}
}

func testArrayString(t *testing.T, result [2]string, expected [2]string) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf(
				"[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayInt32(t *testing.T, result [2]int32, expected [2]int32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf(
				"[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayInt64(t *testing.T, result [2]int64, expected [2]int64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayUint32(t *testing.T, result [2]uint32, expected [2]uint32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayUint64(t *testing.T, result [2]uint64, expected [2]uint64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayFloat32(t *testing.T, result [2]float32, expected [2]float32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testArrayFloat64(t *testing.T, result [2]float64, expected [2]float64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}

func testSliceString(t *testing.T, result []string, expected []string) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf(
				"[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceInt32(t *testing.T, result []int32, expected []int32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf(
				"[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceInt64(t *testing.T, result []int64, expected []int64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceUint32(t *testing.T, result []uint32, expected []uint32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceUint64(t *testing.T, result []uint64, expected []uint64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceFloat32(t *testing.T, result []float32, expected []float32) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
func testSliceFloat64(t *testing.T, result []float64, expected []float64) {
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("[%d] should be '%#v', got '%#v'", i, expected[i], v)
		}
	}
}
