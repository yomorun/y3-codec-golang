package y3

import (
	"errors"

	encoding2 "github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"

	codec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// DecodePrimitivePacket 将一个完整的Packet的buffer全部读入，返回BasePacket对象
//
// Examples:
// [0x01, 0x01, 0x01] -> Key=0x01, Value=0x01
// [0x41, 0x06, 0x03, 0x01, 0x61, 0x04, 0x01, 0x62] -> key=0x03, value=0x61; key=0x04, value=0x62
func DecodePrimitivePacket(buf []byte) (packet *PrimitivePacket, endPos int, sizeL int, err error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")
	logger.Debugf("buf=%#X", buf)

	if buf == nil || len(buf) < primitivePacketBufferMinimalLength {
		return nil, 0, 0, errors.New("invalid y3 packet minimal size")
	}

	p := &PrimitivePacket{valbuf: buf}

	var pos = 0
	// first byte is `Tag`
	p.tag = codec.NewTag(buf[pos])
	pos++

	// read `Varint` from buf for `Length of value`
	tmpBuf := buf[pos:]
	var bufLen int32
	codec := encoding2.VarCodec{}
	err = codec.DecodePVarInt32(tmpBuf, &bufLen)
	if err != nil {
		return nil, 0, 0, err
	}
	sizeL = codec.Size

	if sizeL < 1 {
		return nil, 0, sizeL, errors.New("malformed, size of Length can not smaller than 1")
	}

	// 根据文档表述，p.length指的是value的长度，所以修改为bufLen的值
	//p.length = uint32(len)
	//pos += int(bufLen)
	p.length = uint32(bufLen)
	pos += sizeL

	endPos = pos + int(p.length)

	logger.Debugf(">>> sizeL=%v, length=%v, pos=%v, endPos=%v", sizeL, p.length, pos, endPos)

	p.valbuf = buf[pos:endPos]
	logger.Debugf("valbuf = %#X", p.valbuf)

	return p, endPos, sizeL, nil
}
