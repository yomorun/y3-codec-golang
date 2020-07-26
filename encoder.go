package y3

import (
	"bytes"
	"fmt"

	encoding "github.com/yomorun/yomo-codec-golang/pkg"
)

// Encoder will encode object to Y3 encoding
type encoder struct {
	seqID  int
	valbuf *bytes.Buffer
	isNode bool
	buf    *bytes.Buffer
}

type iEncoder interface {
	Encode() []byte
}

func (enc *encoder) String() string {
	return fmt.Sprintf("Encoder: isNode=%v | seqID=%#x | valbuf=%#v | buf=%#v", enc.isNode, enc.seqID, enc.valbuf, enc.buf)
}

// PirmitivePacketEncoder used for encode a primitive packet
type PirmitivePacketEncoder struct {
	encoder
}

// NewPrimitivePacketEncoder return an encoder for primitive packet
func NewPrimitivePacketEncoder(sid int) *PirmitivePacketEncoder {
	primEnc := &PirmitivePacketEncoder{
		encoder: encoder{
			isNode: false,
			buf:    new(bytes.Buffer),
			valbuf: new(bytes.Buffer),
		},
	}

	primEnc.seqID = sid
	return primEnc
}

// SetInt32Value encode int32 value
func (enc *PirmitivePacketEncoder) SetInt32Value(v int32) {
	size := encoding.SizeOfPVarInt32(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodePVarInt32(buf, v)
	if err != nil {
		panic(err)
	}
	enc.valbuf.Write(buf)
}

// SetStringValue encode string
func (enc *PirmitivePacketEncoder) SetStringValue(v string) {
	buf := []byte(v)
	enc.valbuf.Write(buf)
}

// NodePacketEncoder used for encode a node packet
type NodePacketEncoder struct {
	encoder
}

// NewNodePacketEncoder returns an encoder for node packet
func NewNodePacketEncoder(sid int) *NodePacketEncoder {
	nodeEnc := &NodePacketEncoder{
		encoder: encoder{
			isNode: true,
			buf:    new(bytes.Buffer),
			valbuf: new(bytes.Buffer),
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
func (enc *NodePacketEncoder) AddPrimitivePacket(np *PirmitivePacketEncoder) {
	enc.addRawPacket(np)
}

func (enc *encoder) addRawPacket(en iEncoder) {
	enc.valbuf.Write(en.Encode())
}

// setTag write tag as seqID
func (enc *encoder) writeTag() {
	if enc.seqID < 0 || enc.seqID > 0x7F {
		panic("sid should be in [0..0x7F]")
	}
	if enc.isNode {
		enc.seqID = enc.seqID | 0x80
	}
	enc.buf.WriteByte(byte(enc.seqID))
}

func (enc *encoder) writeLengthBuf() {
	vallen := enc.valbuf.Len()
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
	enc.buf.Write(enc.valbuf.Bytes())

	return enc.buf.Bytes()
}
