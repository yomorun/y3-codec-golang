package tester

type BasicTestData struct {
	Vstring  string  `yomo:"0x10"`
	Vint32   int32   `yomo:"0x11"`
	Vint64   int64   `yomo:"0x12"`
	Vuint32  uint32  `yomo:"0x13"`
	Vuint64  uint64  `yomo:"0x14"`
	Vfloat32 float32 `yomo:"0x15"`
	Vfloat64 float64 `yomo:"0x16"`
	Vbool    bool    `yomo:"0x17"`
}

type EmbeddedTestData struct {
	BasicTestData `yomo:"0x1a"`
	Vaction       string `yomo:"0x1b"`
}

type EmbeddedMoreTestData struct {
	EmbeddedTestData `yomo:"0x1c"`
	Vanimal          string `yomo:"0x1d"`
}

type NamedTestData struct {
	Base    BasicTestData `yomo:"0x1e"`
	Vaction string        `yomo:"0x1f"`
}

type NamedMoreTestData struct {
	MyNest  NamedTestData `yomo:"0x2a"`
	Vanimal string        `yomo:"0x2b"`
}

type ArrayTestData struct {
	Vfoo          string     `yomo:"0x20"`
	Vbar          [2]string  `yomo:"0x21"`
	Vint32Array   [2]int32   `yomo:"0x22"`
	Vint64Array   [2]int64   `yomo:"0x23"`
	Vuint32Array  [2]uint32  `yomo:"0x24"`
	Vuint64Array  [2]uint64  `yomo:"0x25"`
	Vfloat32Array [2]float32 `yomo:"0x26"`
	Vfloat64Array [2]float64 `yomo:"0x27"`
}

type SliceTestData struct {
	Vfoo          string    `yomo:"0x30"`
	Vbar          []string  `yomo:"0x31"`
	Vint32Slice   []int32   `yomo:"0x32"`
	Vint64Slice   []int64   `yomo:"0x33"`
	Vuint32Slice  []uint32  `yomo:"0x34"`
	Vuint64Slice  []uint64  `yomo:"0x35"`
	Vfloat32Slice []float32 `yomo:"0x36"`
	Vfloat64Slice []float64 `yomo:"0x37"`
}

type SliceStructTestData struct {
	Vstring          string                 `yomo:"0x2e"`
	BaseList         []BasicTestData        `yomo:"0x2f"`
	NamedMoreList    []NamedMoreTestData    `yomo:"0x3a"`
	EmbeddedMoreList []EmbeddedMoreTestData `yomo:"0x3b"`
}

type ArrayStructTestData struct {
	Vstring          string                  `yomo:"0x2e"`
	BaseList         [2]BasicTestData        `yomo:"0x2f"`
	NamedMoreList    [2]NamedMoreTestData    `yomo:"0x3a"`
	EmbeddedMoreList [2]EmbeddedMoreTestData `yomo:"0x3b"`
}

type NestedTestData struct {
	SubNested Sub1NestedTestData `yomo:"0x3a"`
}

type Sub1NestedTestData struct {
	SubNested Sub2NestedTestData `yomo:"0x3b"`
}

type Sub2NestedTestData struct {
	SubNested Sub3NestedTestData `yomo:"0x3c"`
}

type Sub3NestedTestData struct {
	BasicList []BasicTestData `yomo:"0x3d"`
}
