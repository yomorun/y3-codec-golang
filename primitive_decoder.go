package y3

import (
	"errors"

	codec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
	encoding "github.com/yomorun/yomo-codec-golang/pkg"
)

// DecodePrimitivePacket 将一个完整的Packet的buffer全部读入，返回BasePacket对象
//
// Examples:
// [0x01, 0x01, 0x01] -> Key=0x01, Value=0x01
func DecodePrimitivePacket(buf []byte) (*PrimitivePacket, int, error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")
	logger.Debugf("buf=%v", buf)

	if buf == nil || len(buf) < primitivePacketBufferMinimalLength {
		return nil, 0, errors.New("invalid y3 packet minimal size")
	}

	p := &PrimitivePacket{valbuf: buf}

	var pos = 0
	// first byte is `Tag`
	p.tag = codec.NewTag(buf[pos])
	pos++

	// read `Varint` from buf as `Length`
	tmpBuf := buf[pos:]
	var bufLen int32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt32(tmpBuf, &bufLen)
	if err != nil {
		return nil, 0, err
	}
	len := codec.Size

	logger.Debugf(">>>len=%v", len)
	if len < 1 {
		return nil, 0, errors.New("malformed, Length can not smaller than 1")
	}
	p.length = uint32(len) // Length的值是Value的字节长度
	pos += int(bufLen)

	// read `Value` raw data, len(raw data) = p.Length - 1
	valLength := p.length
	// p.valbuf = make([]byte, valLength)
	endPos := pos + int(valLength)
	// copied := copy(p.valbuf, buf[pos:uint32(pos)+valLength])
	p.valbuf = buf[pos : uint32(pos)+valLength]
	// logger.Debugf("copied raw data length = %v", copied)

	return p, endPos, nil
}
