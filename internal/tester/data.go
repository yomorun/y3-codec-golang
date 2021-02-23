package tester

type BasicTestData struct {
	Vstring  string  `y3:"0x10"`
	Vint32   int32   `y3:"0x11"`
	Vint64   int64   `y3:"0x12"`
	Vuint32  uint32  `y3:"0x13"`
	Vuint64  uint64  `y3:"0x14"`
	Vfloat32 float32 `y3:"0x15"`
	Vfloat64 float64 `y3:"0x16"`
	Vbool    bool    `y3:"0x17"`
}

type EmbeddedTestData struct {
	BasicTestData `y3:"0x1a"`
	Vaction       string `y3:"0x1b"`
}

type EmbeddedMoreTestData struct {
	EmbeddedTestData `y3:"0x1c"`
	Vanimal          string `y3:"0x1d"`
}

type NamedTestData struct {
	Base    BasicTestData `y3:"0x1e"`
	Vaction string        `y3:"0x1f"`
}

type NamedMoreTestData struct {
	MyNest  NamedTestData `y3:"0x2a"`
	Vanimal string        `y3:"0x2b"`
}

type ArrayTestData struct {
	Vfoo          string     `y3:"0x20"`
	Vbar          [2]string  `y3:"0x21"`
	Vint32Array   [2]int32   `y3:"0x22"`
	Vint64Array   [2]int64   `y3:"0x23"`
	Vuint32Array  [2]uint32  `y3:"0x24"`
	Vuint64Array  [2]uint64  `y3:"0x25"`
	Vfloat32Array [2]float32 `y3:"0x26"`
	Vfloat64Array [2]float64 `y3:"0x27"`
}

type SliceTestData struct {
	Vfoo          string    `y3:"0x30"`
	Vbar          []string  `y3:"0x31"`
	Vint32Slice   []int32   `y3:"0x32"`
	Vint64Slice   []int64   `y3:"0x33"`
	Vuint32Slice  []uint32  `y3:"0x34"`
	Vuint64Slice  []uint64  `y3:"0x35"`
	Vfloat32Slice []float32 `y3:"0x36"`
	Vfloat64Slice []float64 `y3:"0x37"`
}

type SliceStructTestData struct {
	Vstring          string                 `y3:"0x2e"`
	BaseList         []BasicTestData        `y3:"0x2f"`
	NamedMoreList    []NamedMoreTestData    `y3:"0x3a"`
	EmbeddedMoreList []EmbeddedMoreTestData `y3:"0x3b"`
}

type ArrayStructTestData struct {
	Vstring          string                  `y3:"0x2e"`
	BaseList         [2]BasicTestData        `y3:"0x2f"`
	NamedMoreList    [2]NamedMoreTestData    `y3:"0x3a"`
	EmbeddedMoreList [2]EmbeddedMoreTestData `y3:"0x3b"`
}

type NestedTestData struct {
	SubNested Sub1NestedTestData `y3:"0x3a"`
}

type Sub1NestedTestData struct {
	SubNested Sub2NestedTestData `y3:"0x3b"`
}

type Sub2NestedTestData struct {
	SubNested Sub3NestedTestData `y3:"0x3c"`
}

type Sub3NestedTestData struct {
	BasicList []BasicTestData `y3:"0x3d"`
}
