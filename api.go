package y3

import (
	"errors"

	"github.com/yomorun/y3-codec-golang/pkg/spec"
)

// EncodeBool encode bool type data
func EncodeBool(tag int, v bool) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetBool(v)
	return p.Encode()
}

// EncodeUInt encode uint type data
func EncodeUInt(tag int, v uint) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetUInt32(uint32(v))
	return p.Encode()
}

// EncodeInt encode int type data
func EncodeInt(tag int, v int) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetInt32(v)
	return p.Encode()
}

// EncodeUInt64 encode uint64 type data
func EncodeUInt64(tag int, v uint64) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetUInt64(v)
	return p.Encode()
}

// EncodeInt64 encode int type data
func EncodeInt64(tag int, v int64) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetInt64(v)
	return p.Encode()
}

// EncodeFloat32 encode float32 type data
func EncodeFloat32(tag int, v float32) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetFloat32(v)
	return p.Encode()
}

// EncodeFloat64 encode float64 type data
func EncodeFloat64(tag int, v float64) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetFloat64(v)
	return p.Encode()
}

// EncodeString encode UTF-8 string data
func EncodeString(tag int, v string) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.SetUTF8String(v)
	return p.Encode()
}

// EncodeBytes encode raw bytes
func EncodeBytes(tag int, v []byte) ([]byte, error) {
	p, err := spec.NewPacket(uint64(tag))
	if err != nil {
		return nil, err
	}
	p.PutBytes(v)
	return p.Encode()
}

// Marshal TODO wip
func Marshal(tag int, obj interface{}) ([]byte, error) {
	panic(errors.New("NotImplementedError"))
}
