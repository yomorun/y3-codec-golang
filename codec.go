package y3

import (
	"reflect"
)

// Y3Codec encode the user's data according to the Y3 encoding rules
type Y3Codec interface {
	// Marshal encode interface to []byte
	Marshal(input interface{}) ([]byte, error)
}

// NewCodec create a Y3Codec interface
func NewCodec(observe byte) Y3Codec {
	return &y3Codec{
		observe: observe,
	}
}

// y3Codec is implementation of the Y3Codec interface
type y3Codec struct {
	observe byte
}

// Marshal encode interface to []byte
func (c y3Codec) Marshal(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		return newStructEncoder(c.observe, structEncoderOptionRoot(rootToken)).Encode(input)
	}
	return newBasicEncoder(c.observe, basicEncoderOptionRoot(rootToken)).Encode(input)
}

// isStruct determine whether an interface is a structure
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
