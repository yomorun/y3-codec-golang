package y3

import (
	"errors"

	varint "github.com/yomorun/yomo-codec-golang/internal"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// Decode 将一个完整的Packet的buffer全部读入，返回BasePacket对象
//
// Examples:
// [0x01, 0x01, 0x01, 0x01] -> Key=0x01, Value=-1
func Decode(buf []byte) (*BasePacket, error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")

	logger.Debugf("buf=%v", buf)

	if buf == nil || len(buf) < PacketBufferMinimalLength {
		return nil, errors.New("invalid y3 packet minimal size")
	}

	p := &BasePacket{raw: buf}

	var pos = 0
	// first byte is `Tag`
	p.Tag = buf[pos]
	pos++

	// read `Varint` from buf as `Length`
	len, bufLen, err := parseVarint(buf, pos)
	if err != nil {
		return nil, err
	}
	p.Length = len
	pos += bufLen

	// read `Type` of a value
	t, err := parseType(buf[pos])
	if err != nil {
		return nil, err
	}
	p.Type = t
	pos++

	// read `Value` by `Type`
	v, err := parseValue(buf, pos)
	if err != nil {
		return nil, err
	}
	p.Val = v

	return p, nil
}

func parseVarint(b []byte, startPos int) (val int64, len int, err error) {
	dec, len := varint.NewDecoder(b, startPos)
	val, err = dec.Decode()
	if err != nil {
		return 0, 0, err
	}
	return val, len, nil
}

func parseType(b byte) (Type, error) {
	t := Type(b)
	err := t.isValid()
	return t, err
}

func parseValue(b []byte, startPos int) (*Val, error) {
	return &Val{raw: b[startPos:]}, nil
}
