package spec

import (
	"github.com/yomorun/y3-codec-golang/pkg/encoding"
)

// NewPacket create new Packet object
func NewPacket(sid uint64) (*Packet, error) {
	var p = &Packet{}
	tmp, err := getPVarUInt64Buffer(sid)
	p.tagbuf = tmp
	return p, err
}

// NewRawPacket create new Packet object with
func NewRawPacket(rawID []byte) (*Packet, error) {
	// TODO, validate rawID
	var p = &Packet{}
	p.tagbuf = rawID
	return p, nil
}

// SetNil set nil value, means length=0 packet
func (p *Packet) SetNil() {
	p.valbuf = make([]byte, 0)
}

// SetUInt32 set UInt32 value
func (p *Packet) SetUInt32(v uint32) error {
	size := encoding.SizeOfPVarUInt32(v)
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodePVarUInt32(p.valbuf, v)
	if err != nil {
		return err
	}
	return nil
}

// SetInt32 set Int32 value
func (p *Packet) SetInt32(v int) error {
	size := encoding.SizeOfPVarInt32(int32(v))
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodePVarInt32(p.valbuf, int32(v))
	if err != nil {
		return err
	}
	return nil
}

// SetUInt64 set UInt32 value
func (p *Packet) SetUInt64(v uint64) error {
	size := encoding.SizeOfPVarUInt64(v)
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodePVarUInt64(p.valbuf, v)
	if err != nil {
		return err
	}
	return nil
}

// SetInt64 set Int32 value
func (p *Packet) SetInt64(v int64) error {
	size := encoding.SizeOfPVarInt64(int64(v))
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodePVarInt64(p.valbuf, int64(v))
	if err != nil {
		return err
	}
	return nil
}

// SetFloat32 set float32 value
func (p *Packet) SetFloat32(v float32) error {
	size := encoding.SizeOfVarFloat32(v)
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodeVarFloat32(p.valbuf, v)
	if err != nil {
		return err
	}
	return nil
}

// SetFloat64 set float64 value
func (p *Packet) SetFloat64(v float64) error {
	size := encoding.SizeOfVarFloat64(v)
	codec := encoding.VarCodec{Size: size}
	p.valbuf = make([]byte, size)
	// set packet valbuf
	err := codec.EncodeVarFloat64(p.valbuf, v)
	if err != nil {
		return err
	}
	return nil
}

// SetBool set boolean value
func (p *Packet) SetBool(v bool) error {
	var val uint64 = 0
	if v {
		val = 1
	}
	return p.SetUInt64(val)
}

// SetUTF8String set string value
func (p *Packet) SetUTF8String(v string) {
	p.valbuf = []byte(v)
}

// PutBytes append bytes value
func (p *Packet) PutBytes(v []byte) {
	p.valbuf = append(p.valbuf, v...)
}

// AddNode add a child Packet
func (p *Packet) AddNode(child *Packet) (*Packet, error) {
	childBuffer, err := child.Encode()
	if err != nil {
		return nil, err
	}
	p.PutBytes(childBuffer)
	return p, nil
}

// Encode return whole bytes of this packet
func (p *Packet) Encode() ([]byte, error) {
	// if tag buffer is none, read from idTag as PVarUint64 type
	if len(p.tagbuf) < 1 {
		tagbuf, err := getPVarUInt64Buffer(p.idTag)
		if err != nil {
			return nil, err
		}
		p.tagbuf = tagbuf
	}

	// set length buffer
	p.Length = uint64(len(p.valbuf))
	lenbuf, err := getPVarUInt64Buffer(p.Length)
	p.lenbuf = lenbuf
	if err != nil {
		return nil, err
	}
	// [Tag][Length][Value]
	res := append(p.tagbuf, p.lenbuf...)
	res = append(res, p.valbuf...)
	p.buffer = res
	return res, nil
}

func getPVarUInt64Buffer(val uint64) ([]byte, error) {
	size := encoding.SizeOfPVarUInt64(val)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodePVarUInt64(buf, val)
	return buf, err
}
