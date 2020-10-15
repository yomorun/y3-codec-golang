package codes

type protoCodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v *interface{}) error
}
