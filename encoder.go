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
	enc.writeLengthBuf()
	enc.buf.Write(enc.valbuf)

	return enc.buf.Bytes()
}

// CreateNodePacket generate new node
func (enc *Encoder) CreateNodePacket(sid int) *Encoder {
	return CreateEncoder()
}

// CreatePrimitivePacket generate new primitive
func (enc *Encoder) CreatePrimitivePacket(sid int) *Encoder {
	primEnc := &Encoder{
		isNode: false,
		buf:    new(bytes.Buffer),
	}

	primEnc.writeTag(sid)
	return primEnc
}

// AddNodePacket add new node to this node
func (enc *Encoder) AddNodePacket(np *Encoder) {

}

// AddPrimitivePacket add new primitive to this node
func (enc *Encoder) AddPrimitivePacket(np *Encoder) {
	enc.insertBuffer(np.Encode())
}

// SetInt64Value encode int64 value
func (enc *Encoder) SetInt64Value(v int64) {
	buf, length, err := encoding.EncodePvarint(v)
	if err != nil {
		panic(err)
	}
	enc.valbuf = make([]byte, length)
	copy(enc.valbuf, buf)

	enc.vallen = length
}

// SetStringValue encode string
func (enc *Encoder) SetStringValue(v string) {
	buf := []byte(v)
	enc.vallen = len(buf)
	enc.valbuf = make([]byte, len(buf))
	copy(enc.valbuf, buf)
}

// insertBuffer add bufer to node
func (enc *Encoder) insertBuffer(buf []byte) {
	_ = copy(enc.valbuf, buf)
}

// setTag write tag as seqID
func (enc *Encoder) writeTag(sid int) {
	if sid < 0 || sid > 0x7F {
		panic("sid should be in [0..0x7F]")
	}
	enc.seqID = sid
	if enc.isNode {
		enc.seqID = sid | 0x80
	}
	enc.buf.WriteRune(rune(enc.seqID))
}

func (enc *Encoder) writeLengthBuf() {
	if enc.vallen < 1 {
		panic("length must greater than 0")
	}

	buf, _, err := encoding.EncodePvarint(int64(enc.vallen))
	if err != nil {
		panic(err)
	}
	enc.buf.Write(buf)
}
