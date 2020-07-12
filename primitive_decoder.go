package y3

import (
	"errors"

	varint "github.com/yomorun/yomo-codec-golang/internal"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// DecodePrimitivePacket 将一个完整的Packet的buffer全部读入，返回BasePacket对象
//
// Examples:
// [0x01, 0x01, 0x01, 0x01] -> Key=0x01, Value=-1
func DecodePrimitivePacket(buf []byte) (*PrimitivePacket, int, error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")
	logger.Debugf("buf=%v", buf)

	if buf == nil || len(buf) < PrimitivePacketBufferMinimalLength {
		return nil, 0, errors.New("invalid y3 packet minimal size")
	}

	p := &PrimitivePacket{basePacket: basePacket{raw: buf}}

	var pos = 0
	// first byte is `Tag`
	p.Tag = buf[pos]
	pos++

	// read `Varint` from buf as `Length`
	len, bufLen, err := varint.ParseVarintLength(buf, pos)
	if err != nil {
		return nil, 0, err
	}
	if len < 2 {
		return nil, 0, errors.New("malformed, Length can not smaller than 2")
	}
	p.basePacket.length = len - 1 // Length的值是Type+Value的字节长度
	pos += bufLen

	// read `Type` of a value
	t, err := parsePrimitiveType(buf[pos])
	if err != nil {
		return nil, 0, err
	}
	p.Type = t
	pos++

	// read `Value` raw data, len(raw data) = p.Length - 1
	valLength := p.Length()
	p.basePacket.raw = make([]byte, valLength)
	endPos := pos + int(valLength)
	copied := copy(p.basePacket.raw, buf[pos:int64(pos)+valLength])
	logger.Debugf("copied raw data length = %v", copied)

	return p, endPos, nil
}

// get packet.Type and check if is valid type defination
func parsePrimitiveType(b byte) (PrimitiveType, error) {
	t := PrimitiveType(b)
	err := t.isValid()
	return t, err
}
