package codes

import (
	"reflect"

	y3 "github.com/yomorun/y3-codec-golang"

	"github.com/yomorun/y3-codec-golang/pkg/codes/packetstructure"
)

// ProtoCodec: proto codec interface, using for YomoCodec
type ProtoCodec interface {

	// Marshal: Marshal interface to []byte
	Marshal(input interface{}) ([]byte, error)
	// MarshalNoWrapper: Marshal interface to []byte, No Outside Nodes
	MarshalNative(input interface{}) ([]byte, error)

	// UnmarshalStruct: Unmarshal struct to interface
	UnmarshalStruct(data []byte, mold interface{}) error

	// UnmarshalBasic: Unmarshal basic type to interface
	UnmarshalBasic(data []byte, mold *interface{}) error

	// UnmarshalStruct: Unmarshal []byte to interface
	Unmarshal(data []byte, moldInfo *MoldInfo) error

	// UnmarshalByNodePacket: Unmarshal NodePacket to interface
	UnmarshalByNodePacket(node *y3.NodePacket, moldInfo *MoldInfo) error
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

// hold the model data
type MoldInfo struct {
	Mold interface{}
}

func (c *protoCodec) Marshal(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		return packetstructure.EncodeToBytesWith(c.Observe, input)
	}
	//return marshalBasicNative(c.Observe, input)
	return encodeBasic(c.Observe, input)
}

func (c *protoCodec) MarshalNative(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		np, err := packetstructure.Encode(input)
		if err != nil {
			return nil, err
		}
		return np.GetValBuf(), nil
	}
	return marshalBasicNative(c.Observe, input)
}

func (c *protoCodec) UnmarshalStruct(data []byte, mold interface{}) error {
	return c.structDecoder.Unmarshal(data, mold)
}

func (c *protoCodec) unmarshalStructNative(data []byte, mold interface{}) error {
	return c.structDecoder.UnmarshalNative(data, mold)
}

func (c *protoCodec) UnmarshalBasic(data []byte, mold *interface{}) error {
	return c.basicDecoder.Unmarshal(data, mold)
}

func (c *protoCodec) unmarshalBasicNative(data []byte, mold *interface{}) error {
	return c.basicDecoder.UnmarshalNative(data, mold)
}

func (c *protoCodec) isStruct(mold interface{}) bool {
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

func (c *protoCodec) Unmarshal(data []byte, moldInfo *MoldInfo) error {
	if c.isStruct(moldInfo.Mold) {
		err := c.unmarshalStructNative(data, moldInfo.Mold)
		if err != nil {
			return err
		}

	} else {
		err := c.unmarshalBasicNative(data, &moldInfo.Mold)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *protoCodec) unmarshalStructByNodePacket(node *y3.NodePacket, mold interface{}) error {
	return c.structDecoder.UnmarshalByNodePacket(node, mold)
}

func (c *protoCodec) unmarshalBasicByNodePacket(node *y3.NodePacket, mold *interface{}) error {
	return c.basicDecoder.UnmarshalByNodePacket(node, mold)
}

func (c *protoCodec) UnmarshalByNodePacket(node *y3.NodePacket, moldInfo *MoldInfo) error {
	if c.isStruct(moldInfo.Mold) {
		err := c.unmarshalStructByNodePacket(node, moldInfo.Mold)
		if err != nil {
			return err
		}
	} else {
		err := c.unmarshalBasicByNodePacket(node, &moldInfo.Mold)
		if err != nil {
			return err
		}
	}
	return nil
}
