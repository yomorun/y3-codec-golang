package y3

import (
	"errors"

	"github.com/yomorun/y3-codec-golang/pkg/spec/encoding"

	"github.com/yomorun/y3-codec-golang/internal/mark"
)

func parsePayload(b []byte) (endPos int, ifNodePacket bool, np *NodePacket, pp *PrimitivePacket, err error) {
	// fmt.Printf("\t(parsePayload) b=%#x\n", b)
	if len(b) == 0 {
		return 0, false, nil, nil, errors.New("parsePacket params can not be nil")
	}

	pos := 0
	// NodePacket
	if ok := isNodePacket(b[pos]); ok {
		np, endPos, err := DecodeNodePacket(b)
		return endPos, true, np, nil, err
	}

	pp, endPos, _, err = DecodePrimitivePacket(b)
	// fmt.Printf("\t\tb=[%#x], pp:%v\n", b, pp)
	return endPos, false, nil, pp, err
}

// DecodeNodePacket parse out whole buffer to a NodePacket
func DecodeNodePacket(b []byte) (pct *NodePacket, endPos int, err error) {
	// fmt.Println(hex.Dump(b))
	pct = &NodePacket{}

	if len(b) == 0 {
		return pct, 0, nil
	}

	nodeArr := make([]NodePacket, 0)
	primArr := make([]PrimitivePacket, 0)

	pos := 0

	// `Tag`
	tag := mark.NewTag(b[pos])
	pct.basePacket.tag = tag
	// fmt.Printf("pos=%d, n.Tag=%v\n", pos, pct.tag.String())
	pos++

	// `Length`: the type is `varint`
	tmpBuf := b[pos:]
	var vallen int32
	codec := encoding.VarCodec{}
	err = codec.DecodePVarInt32(tmpBuf, &vallen)
	// _len, vallen, err := encoding.Upvarint(b, pos)
	if err != nil {
		return nil, 0, err
	}

	//pct.basePacket.length = uint32(codec.Size) // Length的值是Value的字节长度
	// fmt.Printf("pos=%v, vallen=%v, n.Length=%v\n", pos, vallen, pct.Length())
	// TODO:根据文档表述，pct.basePacket.length指的是value的长度，所以修改为vallen的值
	pct.basePacket.length = uint32(vallen)
	pos += codec.Size

	// `raw` is pct.Length() length
	vl := int(vallen)
	endPos = pos + vl
	pct.basePacket.valbuf = make([]byte, vl)
	copy(pct.basePacket.valbuf, b[pos:endPos])
	// fmt.Printf("pos=%v, endPos=%v, buf=[% x], n.Raw=[% x](len=%v)\n", pos, endPos, b[pos:endPos], pct.basePacket.valbuf, len(pct.basePacket.valbuf))

	// Parse value to Packet
	for {
		// fmt.Println("------>pos:", pos, ", endPos:", endPos, ", len(b)", len(b))
		if pos >= endPos || pos >= len(b) {
			// fmt.Println("===GAME OVER===")
			break
		}
		_p, isNode, np, pp, err := parsePayload(b[pos:endPos])
		pos += _p
		if err != nil {
			return nil, 0, err
		}
		if isNode {
			// fmt.Printf("\tisNode=true, _p=%v, pos=%v, payload=% x, [np]:%v\n", _p, pos, b[pos:endPos], np)
			nodeArr = append(nodeArr, *np)
		} else {
			// fmt.Printf("\tisNode=false, _p=%v, pos=%v, payload=% x, [pp]:%v\n", _p, pos, b[pos:endPos], pp)
			primArr = append(primArr, *pp)
		}
	}

	pct.NodePackets = nodeArr
	pct.PrimitivePackets = primArr

	return pct, endPos, nil
}
