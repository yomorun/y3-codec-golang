package y3

import (
	"reflect"
)

type Y3Codec interface {
	// Marshal: Marshal interface to []byte
	Marshal(input interface{}) ([]byte, error)
}

func NewCodec(observe byte) Y3Codec {
	return &y3Codec{
		observe: observe,
	}
}

type y3Codec struct {
	observe byte
}

func (c y3Codec) Marshal(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		return NewStructEncoderWithRoot(c.observe, input, rootToken).Encode(input)
	}
	return NewBasicEncoderWithRoot(c.observe, rootToken).Encode(input)
}

func (c y3Codec) isStruct(mold interface{}) bool {
	isStruct := false

	moldValue := reflect.Indirect(reflect.ValueOf(mold))
	moldType := moldValue.Type()
	switch moldType.Kind() {
	case reflect.Struct:
		isStruct = true
	case reflect.Slice:
		if moldType.Elem().Kind() == reflect.Struct {
			isStruct = true
		}
	}

	return isStruct
}
