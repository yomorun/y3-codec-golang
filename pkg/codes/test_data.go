package codes

type CodecBenchmarkData struct {
}

func NewCodecBenchmarkData() *CodecBenchmarkData {
	return &CodecBenchmarkData{}
}

func (d *CodecBenchmarkData) DefaultPersonData() []byte {
	input := Person{
		Name: "zhang san",
		Age:  25,
	}

	proto := NewProtoCodec(byte(0))
	inputBuf, _ := proto.Marshal(input)

	return inputBuf
}

func (d *CodecBenchmarkData) DefaultPersonMaxData() []byte {
	person := Person{
		Name: "zhang san",
		Age:  25,
	}
	input := PersonMax{
		X10: uint32(16),
		X11: uint32(17),
		X12: uint32(18),
		X13: uint32(19),
		X14: uint32(20),
		X15: uint32(21),
		X16: uint32(22),
		X17: uint32(23),
		X18: uint32(24),
		X19: uint32(25),
		X1a: uint32(26),
		X1b: uint32(27),
		X1c: uint32(28),
		X1d: uint32(29),
		X1e: uint32(30),
		X1f: uint32(31),
		X20: uint32(32),
		X21: uint32(33),
		X22: uint32(34),
		X23: person,
		X26: uint32(38),
		X27: uint32(39),
		X28: uint32(40),
		X29: uint32(41),
		X2a: uint32(42),
		X2b: uint32(43),
		X2c: uint32(44),
		X2d: uint32(45),
		X2e: uint32(46),
		X2f: uint32(47),
		X30: uint32(48),
		X31: uint32(49),
		X32: uint32(50),
		X33: uint32(51),
		X34: uint32(52),
		X35: uint32(53),
		X36: uint32(54),
		X37: uint32(55),
		X38: uint32(56),
		X39: uint32(57),
		X3a: uint32(58),
		X3b: uint32(59),
		X3c: uint32(60),
		X3d: uint32(61),
		X3e: uint32(62),
		X3f: uint32(63),
	}

	proto := NewProtoCodec(byte(0))
	inputBuf, _ := proto.Marshal(input)

	return inputBuf
}

type Person struct {
	Name string `yomo:"0x24"`
	Age  uint32 `yomo:"0x35"`
}

type PersonMax struct {
	X10 uint32 `yomo:"0x10"`
	X11 uint32 `yomo:"0x11"`
	X12 uint32 `yomo:"0x12"`
	X13 uint32 `yomo:"0x13"`
	X14 uint32 `yomo:"0x14"`
	X15 uint32 `yomo:"0x15"`
	X16 uint32 `yomo:"0x16"`
	X17 uint32 `yomo:"0x17"`
	X18 uint32 `yomo:"0x18"`
	X19 uint32 `yomo:"0x19"`
	X1a uint32 `yomo:"0x1a"`
	X1b uint32 `yomo:"0x1b"`
	X1c uint32 `yomo:"0x1c"`
	X1d uint32 `yomo:"0x1d"`
	X1e uint32 `yomo:"0x1e"`
	X1f uint32 `yomo:"0x1f"`
	X20 uint32 `yomo:"0x20"`
	X21 uint32 `yomo:"0x21"`
	X22 uint32 `yomo:"0x22"`
	X23 Person `yomo:"0x23"`
	//X24 uint32 `yomo:"0x24"`
	//X25 uint32 `yomo:"0x25"`
	X26 uint32 `yomo:"0x26"`
	X27 uint32 `yomo:"0x27"`
	X28 uint32 `yomo:"0x28"`
	X29 uint32 `yomo:"0x29"`
	X2a uint32 `yomo:"0x2a"`
	X2b uint32 `yomo:"0x2b"`
	X2c uint32 `yomo:"0x2c"`
	X2d uint32 `yomo:"0x2d"`
	X2e uint32 `yomo:"0x2e"`
	X2f uint32 `yomo:"0x2f"`
	X30 uint32 `yomo:"0x30"`
	X31 uint32 `yomo:"0x31"`
	X32 uint32 `yomo:"0x32"`
	X33 uint32 `yomo:"0x33"`
	X34 uint32 `yomo:"0x34"`
	X35 uint32 `yomo:"0x35"`
	X36 uint32 `yomo:"0x36"`
	X37 uint32 `yomo:"0x37"`
	X38 uint32 `yomo:"0x38"`
	X39 uint32 `yomo:"0x39"`
	X3a uint32 `yomo:"0x3a"`
	X3b uint32 `yomo:"0x3b"`
	X3c uint32 `yomo:"0x3c"`
	X3d uint32 `yomo:"0x3d"`
	X3e uint32 `yomo:"0x3e"`
	X3f uint32 `yomo:"0x3f"`
}
