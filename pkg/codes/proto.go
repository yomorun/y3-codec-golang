package codes

import (
	"reflect"

	y3 "github.com/yomorun/yomo-codec-golang"

	"github.com/yomorun/yomo-codec-golang/pkg/codes/packetstructure"
)

// ProtoCodec: proto codec interface, using for YomoCodec
type ProtoCodec interface {

	// Marshal: Marshal interface to []byte
	Marshal(input interface{}) ([]byte, error)
	// MarshalNoWrapper: Marshal interface to []byte, No Outside Nodes
	MarshalNative(input interface{}) ([]byte, error)

	// UnmarshalStruct: Unmarshal struct to interface
	UnmarshalStruct(data []byte, mold interface{}) error
	// UnmarshalStructNative: Unmarshal struct to interface by native data, No Outside Nodes
	UnmarshalStructNative(data []byte, mold interface{}) error

	// UnmarshalBasic: Unmarshal basic type to interface
	UnmarshalBasic(data []byte, mold *interface{}) error
	// UnmarshalBasicNative: Unmarshal basic type to interface by native data, No Outside Nodes
	UnmarshalBasicNative(data []byte, mold *interface{}) error

	// IsStruct: mold is Struct?
	IsStruct(mold interface{}) bool

	// @deprecated
	UnmarshalStructByNodePacket(node *y3.NodePacket, mold interface{}) error
	// @deprecated
	UnmarshalBasicByNodePacket(node *y3.NodePacket, mold *interface{}) error
}

// protoCodec: Implementation of the ProtoCodec Interface
type protoCodec struct {
	Observe       byte
	basicDecoder  *BasicDecoder
	structDecoder *StructDecoder
}

func NewProtoCodec(observe byte) ProtoCodec {
	return &protoCodec{
		Observe:       observe,
		basicDecoder:  newBasicDecoder(observe),
		structDecoder: newStructDecoder(observe),
	}
}

func (c *protoCodec) Marshal(input interface{}) ([]byte, error) {
	if c.IsStruct(input) {
		return packetstructure.EncodeToBytesWith(c.Observe, input)
	}
	//return marshalBasicNative(c.Observe, input)
	return encodeBasic(c.Observe, input)
}

func (c *protoCodec) MarshalNative(input interface{}) ([]byte, error) {
	if c.IsStruct(input) {
		np, err := packetstructure.Encode(input)
		if err != nil {
			return nil, err
		}
		return np.GetValbuf(), nil
	}
	return marshalBasicNative(c.Observe, input)
}

func (c *protoCodec) UnmarshalStruct(data []byte, mold interface{}) error {
	return c.structDecoder.Unmarshal(data, mold)
}

func (c *protoCodec) UnmarshalStructNative(data []byte, mold interface{}) error {
	return c.structDecoder.UnmarshalNative(data, mold)
}

func (c *protoCodec) UnmarshalBasic(data []byte, mold *interface{}) error {
	return c.basicDecoder.Unmarshal(data, mold)
}

func (c *protoCodec) UnmarshalBasicNative(data []byte, mold *interface{}) error {
	return c.basicDecoder.UnmarshalNative(data, mold)
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

// @deprecated
func (c *protoCodec) UnmarshalStructByNodePacket(node *y3.NodePacket, mold interface{}) error {
	return c.structDecoder.UnmarshalByNodePacket(node, mold)
}

// @deprecated
func (c *protoCodec) UnmarshalBasicByNodePacket(node *y3.NodePacket, mold *interface{}) error {
	return c.basicDecoder.UnmarshalByNodePacket(node, mold)
}
