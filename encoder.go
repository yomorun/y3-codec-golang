package y3

import (
	"bytes"
	"fmt"

	"github.com/yomorun/y3-codec-golang/pkg/encoding"
)

// Encoder will encode object to Y3 encoding
type encoder struct {
	seqID   int
	valbuf  []byte
	isNode  bool
	isArray bool
	buf     *bytes.Buffer
}

type iEncoder interface {
	Encode() []byte
}

func (enc *encoder) GetValBuf() []byte {
	return enc.valbuf
}

func (enc *encoder) IsEmpty() bool {
	return len(enc.valbuf) == 0
}

func (enc *encoder) String() string {
	return fmt.Sprintf("Encoder: isNode=%v | seqID=%#x | valBuf=%#v | buf=%#v", enc.isNode, enc.seqID, enc.valbuf, enc.buf)
}

// PrimitivePacketEncoder used for encode a primitive packet
type PrimitivePacketEncoder struct {
	encoder
}

// NewPrimitivePacketEncoder return an Encoder for primitive packet
func NewPrimitivePacketEncoder(sid int) *PrimitivePacketEncoder {
	primEnc := &PrimitivePacketEncoder{
		encoder: encoder{
			isNode: false,
			buf:    new(bytes.Buffer),
		},
	}

	primEnc.seqID = sid
	return primEnc
}

// SetInt32Value encode int32 value
func (enc *PrimitivePacketEncoder) SetInt32Value(v int32) {
	size := encoding.SizeOfPVarInt32(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodePVarInt32(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
	// enc.valBuf.Write(buf)
}

// SetUInt32Value encode uint32 value
func (enc *PrimitivePacketEncoder) SetUInt32Value(v uint32) {
	size := encoding.SizeOfPVarUInt32(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodePVarUInt32(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetInt64Value encode int64 value
func (enc *PrimitivePacketEncoder) SetInt64Value(v int64) {
	size := encoding.SizeOfPVarInt64(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodePVarInt64(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetUInt64Value encode uint64 value
func (enc *PrimitivePacketEncoder) SetUInt64Value(v uint64) {
	size := encoding.SizeOfPVarUInt64(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodePVarUInt64(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetFloat32Value encode float32 value
func (enc *PrimitivePacketEncoder) SetFloat32Value(v float32) {
	var size = encoding.SizeOfVarFloat32(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodeVarFloat32(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetFloat64Value encode float64 value
func (enc *PrimitivePacketEncoder) SetFloat64Value(v float64) {
	var size = encoding.SizeOfVarFloat64(v)
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodeVarFloat64(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetBoolValue encode bool value
func (enc *PrimitivePacketEncoder) SetBoolValue(v bool) {
	var size = encoding.SizeOfPVarUInt32(uint32(1))
	codec := encoding.VarCodec{Size: size}
	enc.valbuf = make([]byte, size)
	err := codec.EncodePVarBool(enc.valbuf, v)
	if err != nil {
		panic(err)
	}
}

// SetStringValue encode string
func (enc *PrimitivePacketEncoder) SetStringValue(v string) {
	// buf := []byte(v)
	// enc.valBuf.Write(buf)
	enc.valbuf = []byte(v)
}

// SetBytes set bytes to internal buf variable
func (enc *PrimitivePacketEncoder) SetBytes(buf []byte) {
	enc.valbuf = buf
}

// NodePacketEncoder used for encode a node packet
type NodePacketEncoder struct {
	encoder
}

// NewNodePacketEncoder returns an Encoder for node packet
func NewNodePacketEncoder(sid int) *NodePacketEncoder {
	nodeEnc := &NodePacketEncoder{
		encoder: encoder{
			isNode: true,
			buf:    new(bytes.Buffer),
		},
	}

	nodeEnc.seqID = sid
	return nodeEnc
}

// NewNodeSlicePacketEncoder returns an Encoder for node packet that is a slice
func NewNodeSlicePacketEncoder(sid int) *NodePacketEncoder {
	nodeEnc := &NodePacketEncoder{
		encoder: encoder{
			isNode:  true,
			isArray: true,
			buf:     new(bytes.Buffer),
		},
	}

	nodeEnc.seqID = sid
	return nodeEnc
}

// AddNodePacket add new node to this node
func (enc *NodePacketEncoder) AddNodePacket(np *NodePacketEncoder) {
	enc.addRawPacket(np)
}

// AddPrimitivePacket add new primitive to this node
func (enc *NodePacketEncoder) AddPrimitivePacket(np *PrimitivePacketEncoder) {
	enc.addRawPacket(np)
}

func (enc *encoder) addRawPacket(en iEncoder) {
	enc.valbuf = append(enc.valbuf, en.Encode()...)
}

func (enc *encoder) AddBytes(buf []byte) {
	enc.valbuf = append(enc.valbuf, buf...)
}

// setTag write tag as seqID
func (enc *encoder) writeTag() {
	//fmt.Printf("#60 enc.seqID=%#x\n", enc.seqID)
	if enc.seqID < 0 || enc.seqID > 0x7F {
		panic("sid should be in [0..0x7F]")
	}
	if enc.isNode {
		enc.seqID = enc.seqID | 0x80
	}
	if enc.isArray {
		enc.seqID = enc.seqID | 0x40
	}
	enc.buf.WriteByte(byte(enc.seqID))
}

func (enc *encoder) writeLengthBuf() {
	// vallen := enc.valBuf.Len()
	vallen := len(enc.valbuf)
	if vallen < 1 {
		panic("length must greater than 0")
	}

	size := encoding.SizeOfPVarInt32(int32(vallen))
	codec := encoding.VarCodec{Size: size}
	tmp := make([]byte, size)
	err := codec.EncodePVarInt32(tmp, int32(vallen))
	// buf, _, err := encoding.EncodePvarint(v)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("tmp=%#x, vallen=%v", tmp, vallen)
	enc.buf.Write(tmp)
}

// Encode returns a final Y3 encoded byte slice
func (enc *encoder) Encode() []byte {
	// Tag
	enc.writeTag()
	// Length
	enc.writeLengthBuf()
	// Value
	// enc.buf.Write(enc.valBuf.Bytes())
	enc.buf.Write(enc.valbuf)

	return enc.buf.Bytes()
}
