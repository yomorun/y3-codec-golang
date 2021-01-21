package encoding

// EncodePVarBool encode value as PVarBool to buffer
func (codec *VarCodec) EncodePVarBool(buffer []byte, value bool) error {
	tmp := int64(1)
	if !value {
		tmp = int64(0)
	}
	return codec.encodePVarInt(buffer, tmp)
}

// DecodePVarBool decode to value as PVarBool from buffer
func (codec *VarCodec) DecodePVarBool(buffer []byte, value *bool) error {
	var tmp int64
	var err = codec.decodePVarInt(buffer, &tmp)
	if tmp == 1 {
		*value = true
	} else {
		*value = false
	}
	return err
}
