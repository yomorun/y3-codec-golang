package y3

import (
	"reflect"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// Codec encode the user's data according to the Y3 encoding rules
type Codec interface {
	// Marshal encode interface to []byte
	Marshal(input interface{}) ([]byte, error)
}

// NewCodec create a Codec interface
func NewCodec(observe byte) Codec {
	return &y3Codec{
		observe: observe,
	}
}

// y3Codec is implementation of the Codec interface
type y3Codec struct {
	observe byte
}

// Marshal encode interface to []byte
func (c y3Codec) Marshal(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		return newStructEncoder(c.observe,
			structEncoderOptionRoot(utils.RootToken),
			structEncoderOptionForbidUserKey(utils.ForbidUserKey),
			structEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
			Encode(input)
	}
	return newBasicEncoder(c.observe,
		basicEncoderOptionRoot(utils.RootToken),
		basicEncoderOptionForbidUserKey(utils.ForbidUserKey),
		basicEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
		Encode(input)
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
