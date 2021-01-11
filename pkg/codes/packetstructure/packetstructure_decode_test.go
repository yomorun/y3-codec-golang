package packetstructure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yomorun/y3-codec-golang/pkg/packetutils"

	y3 "github.com/yomorun/y3-codec-golang"
)

type Basic struct {
	Vstring  string  `yomo:"0x10"`
	Vint32   int32   `yomo:"0x11"`
	Vint64   int64   `yomo:"0x12"`
	Vuint32  uint32  `yomo:"0x13"`
	Vuint64  uint64  `yomo:"0x14"`
	Vfloat32 float32 `yomo:"0x15"`
	Vfloat64 float64 `yomo:"0x16"`
	Vbool    bool    `yomo:"0x17"`
}

func newBasic() Basic {
	return Basic{
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

type Embedded struct {
	Basic   `yomo:"0x1a"`
	Vaction string `yomo:"0x1b"`
}

type EmbeddedMore struct {
	Embedded `yomo:"0x1c"`
	Vanimal  string `yomo:"0x1d"`
}

type Named struct {
	Base    Basic  `yomo:"0x1e"`
	Vaction string `yomo:"0x1f"`
}

type NamedMore struct {
	MyNest  Named  `yomo:"0x2a"`
	Vanimal string `yomo:"0x2b"`
}

type Array struct {
	Vfoo          string     `yomo:"0x20"`
	Vbar          [2]string  `yomo:"0x21"`
	Vint32Array   [2]int32   `yomo:"0x22"`
	Vint64Array   [2]int64   `yomo:"0x23"`
	Vuint32Array  [2]uint32  `yomo:"0x24"`
	Vuint64Array  [2]uint64  `yomo:"0x25"`
	Vfloat32Array [2]float32 `yomo:"0x26"`
	Vfloat64Array [2]float64 `yomo:"0x27"`
}

type Slice struct {
	Vfoo          string    `yomo:"0x30"`
	Vbar          []string  `yomo:"0x31"`
	Vint32Slice   []int32   `yomo:"0x32"`
	Vint64Slice   []int64   `yomo:"0x33"`
	Vuint32Slice  []uint32  `yomo:"0x34"`
	Vuint64Slice  []uint64  `yomo:"0x35"`
	Vfloat32Slice []float32 `yomo:"0x36"`
	Vfloat64Slice []float64 `yomo:"0x37"`
}

type SliceStruct struct {
	Vstring          string         `yomo:"0x2e"`
	BaseList         []Basic        `yomo:"0x2f"`
	NamedMoreList    []NamedMore    `yomo:"0x3a"`
	EmbeddedMoreList []EmbeddedMore `yomo:"0x3b"`
}

type ArrayStruct struct {
	Vstring          string          `yomo:"0x2e"`
	BaseList         [2]Basic        `yomo:"0x2f"`
	NamedMoreList    [2]NamedMore    `yomo:"0x3a"`
	EmbeddedMoreList [2]EmbeddedMore `yomo:"0x3b"`
}

type Nested struct {
	SubNested Sub1Nested `yomo:"0x3a"`
}

type Sub1Nested struct {
	SubNested Sub2Nested `yomo:"0x3b"`
}

type Sub2Nested struct {
	SubNested Sub3Nested `yomo:"0x3c"`
}

type Sub3Nested struct {
	BasicList []Basic `yomo:"0x3d"`
}

func TestBasic_Struct(t *testing.T) {
	t.Parallel()

	input := newBasic()
	node, _ := Encode(input)

	//debug:
	fmt.Println("#404 DEBUG::TestBasic_Struct:")
	fmt.Println(input)
	packetutils.PrintNodePacket(node)
	fmt.Println()

	//result := Basic{}
	var result Basic
	runDecode(t, node, &result)
	testBasicStruct(t, result, input)
	//fmt.Printf("#41.2 result=%v\n", result)
}

func runDecode(t *testing.T, node *y3.NodePacket, output interface{}) {
	err := Decode(node, output)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}
}

func TestDecode_Embedded(t *testing.T) {
	t.Parallel()

	input := Embedded{
		Basic:   newBasic(),
		Vaction: "drink",
	}
	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestDecode_Embedded: ")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result Embedded
	err := Decode(node, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	if result.Vaction != "drink" {
		t.Errorf("vstring value should be 'drink': %#v", result.Vaction)
	}

	testBasicStruct(t, result.Basic, input.Basic)
}

func TestDecode_EmbeddedMore(t *testing.T) {
	t.Parallel()

	input := EmbeddedMore{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"}
	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestDecode_EmbeddedMore: ")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result EmbeddedMore
	err := Decode(node, &result)
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

	testBasicStruct(t, result.Basic, input.Basic)
}

func TestDecoder_Named(t *testing.T) {
	t.Parallel()

	input := Named{Base: newBasic(), Vaction: "drink"}
	node, _ := Encode(input)

	// debug:
	//fmt.Println("#404 DEBUG::TestDecoder_Named: ")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result Named
	err := Decode(node, &result)
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

	input := NamedMore{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}
	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestDecode_EmbeddedMore: ")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result NamedMore
	err := Decode(node, &result)
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

	input := Array{
		"foo",
		[2]string{"foo", "bar"},
		[2]int32{1, 2},
		[2]int64{1, 2},
		[2]uint32{1, 2},
		[2]uint64{1, 2},
		[2]float32{1, 2},
		[2]float64{1, 2},
	}
	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestArray:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result Array
	err := Decode(node, &result)
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

	input := Slice{
		"foo",
		[]string{"foo", "bar"},
		[]int32{1, 2},
		[]int64{1, 2},
		[]uint32{1, 2},
		[]uint64{1, 2},
		[]float32{1, 2},
		[]float64{1, 2},
	}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestSlice:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result Slice
	err := Decode(node, &result)
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

	input := SliceStruct{
		Vstring:  "foo",
		BaseList: []Basic{newBasic(), newBasic()},
		NamedMoreList: []NamedMore{
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: []EmbeddedMore{
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestSliceStruct:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result SliceStruct
	err := Decode(node, &result)
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
		testBasicStruct(t, v.Basic, input.EmbeddedMoreList[i].Basic)
	}
}

func TestArrayStruct(t *testing.T) {
	t.Parallel()

	input := ArrayStruct{
		Vstring:  "foo",
		BaseList: [2]Basic{newBasic(), newBasic()},
		NamedMoreList: [2]NamedMore{
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: [2]EmbeddedMore{
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestArrayStruct:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result ArrayStruct
	err := Decode(node, &result)
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
		testBasicStruct(t, v.Basic, input.EmbeddedMoreList[i].Basic)
	}
}

func TestRootSliceWithBasicStruct(t *testing.T) {
	t.Parallel()

	input := []Basic{newBasic(), newBasic()}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestRootSliceWithBasicStruct:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result []Basic
	err := Decode(node, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValueWith(t, len(result), len(input), "[]Basic len")

	for i, v := range result {
		testBasicStruct(t, v, input[i])
	}
}

func TestRootSliceWithSliceStruct(t *testing.T) {
	t.Parallel()

	input1 := SliceStruct{
		Vstring:  "foo",
		BaseList: []Basic{newBasic(), newBasic()},
		NamedMoreList: []NamedMore{
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{MyNest: Named{Base: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
		EmbeddedMoreList: []EmbeddedMore{
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"},
			{Embedded: Embedded{Basic: newBasic(), Vaction: "drink"}, Vanimal: "bird"}},
	}

	input := []SliceStruct{input1, input1}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestRootSliceWithSliceStruct:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result []SliceStruct
	err := Decode(node, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testValueWith(t, len(result), len(input), "[]SliceStruct len")

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
			testBasicStruct(t, vv.Basic, input[i].EmbeddedMoreList[i].Basic)
		}
	}

}

func TestNested(t *testing.T) {
	t.Parallel()

	input := Nested{Sub1Nested{Sub2Nested{Sub3Nested{
		BasicList: []Basic{newBasic(), newBasic()},
	}}}}

	node, _ := Encode(input)

	//debug:
	//fmt.Println("#404 DEBUG::TestNested:")
	//codes.PrintNodePacket(node)
	//fmt.Println()

	var result Nested
	err := Decode(node, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	for i, v := range result.SubNested.SubNested.SubNested.BasicList {
		testBasicStruct(t, v, input.SubNested.SubNested.SubNested.BasicList[i])
	}

}

func testBasicStruct(t *testing.T, result Basic, expected Basic) {
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
