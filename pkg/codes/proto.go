package codes

type protoCodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v *interface{}) error
	//Unmarshal2(data []byte, v interface{}) error
	//Unmarshal3(data []byte, mold interface{}) (interface{}, error)

	UnmarshalStruct(data []byte, v interface{}) error
	UnmarshalBasic(data []byte, v *interface{}) error
}
