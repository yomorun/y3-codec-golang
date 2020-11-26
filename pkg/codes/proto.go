package codes

import (
	"reflect"

	"github.com/yomorun/yomo-codec-golang/pkg/codes/packetstructure"
)

// ProtoCodec: proto codec interface, using for YomoCodec
type ProtoCodec interface {
	// Marshal: Marshal interface to []byte
	Marshal(input interface{}) ([]byte, error)
	// MarshalNoWrapper: Marshal interface to []byte, but no Outside Nodes
	MarshalNoWrapper(input interface{}) ([]byte, error)
	// UnmarshalStruct: Unmarshal struct to interface
	UnmarshalStruct(data []byte, mold interface{}) error
	// UnmarshalBasic: Unmarshal basic type to interface
	UnmarshalBasic(data []byte, mold *interface{}) error
	// IsStruct: mold is Struct?
	IsStruct(mold interface{}) bool
}

// protoCodec: Implementation of the ProtoCodec Interface
type protoCodec struct {
	Observe string
}

// NewProtoCodec: create a ProtoCodec interface
func NewProtoCodec(observe string) ProtoCodec {
	return &protoCodec{Observe: observe}
}

func (c *protoCodec) Marshal(input interface{}) ([]byte, error) {
	if c.IsStruct(input) {
		return packetstructure.EncodeToBytesWith(c.Observe, input)
	}
	return marshalPrimitive(c.Observe, input)
}

func (c *protoCodec) MarshalNoWrapper(input interface{}) ([]byte, error) {
	if c.IsStruct(input) {
		np, err := packetstructure.Encode(input)
		if err != nil {
			return nil, err
		}
		return np.GetValbuf(), nil
	}
	return marshalPrimitive(c.Observe, input)
}

func (c *protoCodec) UnmarshalStruct(data []byte, mold interface{}) error {
	decoder := newStructDecoder(c.Observe)
	return decoder.Unmarshal(data, mold)
}

func (c *protoCodec) UnmarshalBasic(data []byte, mold *interface{}) error {
	decoder := newBasicDecoder(c.Observe)
	return decoder.Unmarshal(data, mold)
}

func (c *protoCodec) IsStruct(mold interface{}) bool {
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
