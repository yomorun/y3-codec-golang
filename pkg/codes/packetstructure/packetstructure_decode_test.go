package packetstructure

import (
	"testing"
)

type Basic struct {
	Vstring  string  `yomo:"0x10"`
	Vint32   int32   `yomo:"0x11"`
	Vint64   int64   `yomo:"0x12"`
	Vuint32  uint32  `yomo:"0x13"`
	Vuint64  uint64  `yomo:"0x14"`
	Vfloat32 float32 `yomo:"0x15"`
	Vfloat64 float64 `yomo:"0x16"`
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
	}
}

type Embedded struct {
	Basic   `yomo:"0x1a"`
	Vaction string `yomo:"0x1b"`
}

type EmbeddedMore struct {
	Embedded `yomo:"0x20"`
	Vanimal  string `yomo:"0x21"`
}

type Named struct {
	Base    Basic  `yomo:"0x1a"`
	Vaction string `yomo:"0x1b"`
}

type NamedMore struct {
	MyNest  Named  `yomo:"0x20"`
	Vanimal string `yomo:"0x21"`
}

type Array struct {
	Vfoo          string     `yomo:"0x25"`
	Vbar          [2]string  `yomo:"0x26"`
	Vint32Array   [2]int32   `yomo:"0x27"`
	Vint64Array   [2]int64   `yomo:"0x28"`
	Vuint32Array  [2]uint32  `yomo:"0x29"`
	Vuint64Array  [2]uint64  `yomo:"0x2a"`
	Vfloat32Array [2]float32 `yomo:"0x2b"`
	Vfloat64Array [2]float64 `yomo:"0x2c"`
}

type Slice struct {
	Vfoo          string    `yomo:"0x25"`
	Vbar          []string  `yomo:"0x26"`
	Vint32Slice   []int32   `yomo:"0x27"`
	Vint64Slice   []int64   `yomo:"0x28"`
	Vuint32Slice  []uint32  `yomo:"0x29"`
	Vuint64Slice  []uint64  `yomo:"0x2a"`
	Vfloat32Slice []float32 `yomo:"0x2b"`
	Vfloat64Slice []float64 `yomo:"0x2c"`
}

func TestBasic_Struct(t *testing.T) {
	t.Parallel()

	input := newBasic()
	node, _ := Encode(input)

	var result Basic
	err := Decode(node, &result)
	if err != nil {
		t.Errorf("got an err: %s", err.Error())
		t.FailNow()
	}

	testBasicStruct(t, result)
}

func TestDecode_Embedded(t *testing.T) {
	t.Parallel()

	input := Embedded{Basic: newBasic(), Vaction: "drink"}
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

	testBasicStruct(t, result.Basic)
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

	testBasicStruct(t, result.Basic)
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

	testBasicStruct(t, result.Base)
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

	testBasicStruct(t, result.MyNest.Base)
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

func testBasicStruct(t *testing.T, value Basic) {
	if value.Vstring != "foo" {
		t.Errorf("vstring value should be 'foo': %#v", value.Vstring)
	}

	if value.Vint32 != 127 {
		t.Errorf("Vint32 value should be 127: %#v", value.Vint32)
	}

	if value.Vint64 != -1 {
		t.Errorf("Vint64 value should be -1: %#v", value.Vint64)
	}

	if value.Vuint32 != 130 {
		t.Errorf("Vuint32 value should be 130: %#v", value.Vuint32)
	}

	if value.Vuint64 != 18446744073709551615 {
		t.Errorf("Vuint64 value should be 18446744073709551615: %#v", value.Vuint64)
	}

	if value.Vfloat32 != 0.25 {
		t.Errorf("Vuint64 value should be 0.25: %#v", value.Vfloat32)
	}

	if value.Vfloat64 != 23 {
		t.Errorf("Vfloat64 value should be 23: %#v", value.Vfloat64)
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
