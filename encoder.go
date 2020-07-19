package y3

import (
	"bytes"

	encoding "github.com/yomorun/yomo-codec-golang/pkg"
)

// Encoder will encode object to Y3 encoding
type Encoder struct {
	seqID  int
	valbuf []byte
	lenbuf []byte
	vallen int
	isNode bool
	buf    *bytes.Buffer
}

// CreateEncoder returns an Encoder
func CreateEncoder() *Encoder {
	return &Encoder{valbuf: make([]byte, 10)}
}

// Encode returns a final Y3 encoded byte slice
func (enc *Encoder) Encode() []byte {
	enc.writeTag()
	enc.writeLengthBuf()
	enc.buf.Write(enc.valbuf)

	return enc.buf.Bytes()
}

// CreateNodePacket generate new node
func (enc *Encoder) CreateNodePacket(sid int) *Encoder {
	nodeEnc := &Encoder{
		isNode: true,
	}
	return nodeEnc
}

// CreatePrimitivePacket generate new primitive
func (enc *Encoder) CreatePrimitivePacket(sid int) *Encoder {
	primEnc := &Encoder{
		isNode: false,
		buf:    new(bytes.Buffer),
	}

	primEnc.seqID = sid
	return primEnc
}

// AddNodePacket add new node to this node
func (enc *Encoder) AddNodePacket(np *Encoder) {

}

// AddPrimitivePacket add new primitive to this node
func (enc *Encoder) AddPrimitivePacket(np *Encoder) {
	enc.valbuf = np.Encode()
	enc.vallen = len(enc.valbuf)
}

// SetInt64Value encode int64 value
func (enc *Encoder) SetInt64Value(v int64) {
	buf, length, err := encoding.EncodePvarint(v)
	if err != nil {
		panic(err)
	}
	enc.valbuf = make([]byte, length)
	copy(enc.valbuf, buf)
}

// SetStringValue encode string
func (enc *Encoder) SetStringValue(v string) {
	buf := []byte(v)
	enc.valbuf = make([]byte, len(buf))
	copy(enc.valbuf, buf)
}

// setTag write tag as seqID
func (enc *Encoder) writeTag() {
	if enc.seqID < 0 || enc.seqID > 0x7F {
		panic("sid should be in [0..0x7F]")
	}
	if enc.isNode {
		enc.seqID = enc.seqID | 0x80
	}
	enc.buf.WriteRune(rune(enc.seqID))
}

func (enc *Encoder) writeLengthBuf() {
	enc.vallen = len(enc.valbuf)
	if enc.vallen < 1 {
		panic("length must greater than 0")
	}

	buf, _, err := encoding.EncodePvarint(int64(enc.vallen))
	if err != nil {
		panic(err)
	}
	enc.buf.Write(buf)
}
