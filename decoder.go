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
func Decode(buf []byte) (*PrimitivePacket, error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")
	logger.Debugf("buf=%v", buf)

	if buf == nil || len(buf) < PrimitivePacketBufferMinimalLength {
		return nil, errors.New("invalid y3 packet minimal size")
	}

	p := &PrimitivePacket{raw: buf}

	var pos = 0
	// first byte is `Tag`
	p.Tag = buf[pos]
	pos++

	// read `Varint` from buf as `Length`
	len, bufLen, err := ParseVarintLength(buf, pos)
	if err != nil {
		return nil, err
	}
	if len < 2 {
		return nil, errors.New("malformed, Length can not smaller than 2")
	}
	p.Length = len - 1 // Length的值是Type+Value的字节长度
	pos += bufLen

	// read `Type` of a value
	t, err := parseType(buf[pos])
	if err != nil {
		return nil, err
	}
	p.Type = t
	pos++

	// read `Value` raw data, len(raw data) = p.Length - 1
	valLength := p.Length
	p.raw = make([]byte, valLength)
	copied := copy(p.raw, buf[pos:int64(pos)+valLength])
	logger.Debugf("copied raw data length = %v", copied)

	return p, nil
}

// ParseVarintLength parse length as a varint type
func ParseVarintLength(b []byte, startPos int) (val int64, len int, err error) {
	dec, len := varint.NewDecoder(b[startPos:])
	val, err = dec.Decode()
	if err != nil {
		return 0, 0, err
	}
	return val, len, nil
}

// get packet.Type and check if is valid type defination
func parseType(b byte) (PrimitiveType, error) {
	t := PrimitiveType(b)
	err := t.isValid()
	return t, err
}
