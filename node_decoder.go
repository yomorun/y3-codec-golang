package y3

import (
	"encoding/hex"
	"errors"
	"fmt"

	varint "github.com/yomorun/yomo-codec-golang/internal"
)

func parsePayload(b []byte) (endPos int, isNodePacket bool, np *NodePacket, pp *PrimitivePacket, err error) {
	fmt.Printf("\t(parsePayload) b=% x\n", b)
	if len(b) == 0 {
		return 0, false, nil, nil, errors.New("parsePacket params can not be nil")
	}

	pos := 0
	// NodePacket
	if ok := IsNodePacket(b[pos]); ok {
		np, endPos, err := ReadNode(b)
		return endPos, true, np, nil, err
	}

	pp, endPos, err = parsePrimitivePacket(b)
	fmt.Printf("\t\tb=[% x], pp:%v\n", b, pp)
	return endPos, false, nil, pp, err
}

func parsePrimitivePacket(b []byte) (pp *PrimitivePacket, endPos int, err error) {
	return DecodePrimitivePacket(b)
}

// ReadNode parse out whole buffer to a NodePacket
func ReadNode(b []byte) (pct *NodePacket, endPos int, err error) {
	fmt.Println(hex.Dump(b))
	pct = &NodePacket{}

	if len(b) == 0 {
		return pct, 0, nil
	}

	nodeArr := make([]NodePacket, 0)
	primArr := make([]PrimitivePacket, 0)

	pos := 0
	// total := len(b)

	// `Tag`
	tag, err := newNodeTag(b[pos])
	if err != nil {
		return nil, 0, err
	}
	pct.Tag = *tag
	fmt.Printf("pos=%d, n.Tag=%v\n", pos, pct.Tag)
	pos++

	// `Length`: the type is `varint`
	_len, lengthOfLenBuffer, err := varint.ParseVarintLength(b, pos)
	if err != nil {
		return nil, 0, err
	}
	pct.basePacket.length = _len // Length的值是Value的字节长度
	fmt.Printf("pos=%v, lengthOfLenBuffer=%v, n.Length=%v\n", pos, lengthOfLenBuffer, pct.Length())
	pos += lengthOfLenBuffer

	// `raw` is pct.Length() length
	vl := int(_len)
	endPos = pos + vl
	pct.basePacket.raw = make([]byte, vl)
	copy(pct.basePacket.raw, b[pos:endPos])
	fmt.Printf("pos=%v, endPos=%v, buf=[% x], n.Raw=[% x](len=%v)\n", pos, endPos, b[pos:endPos], pct.basePacket.raw, len(pct.basePacket.raw))

	// Parse value to Packet
	for {
		fmt.Println("------>pos:", pos, ", endPos:", endPos, ", len(b)", len(b))
		if pos >= endPos || pos >= len(b) {
			fmt.Println("===GAME OVER===")
			break
		}
		_p, isNode, np, pp, err := parsePayload(b[pos:endPos])
		pos += _p
		if err != nil {
			return nil, 0, err
		}
		if isNode {
			fmt.Printf("\tisNode=true, _p=%v, pos=%v, payload=% x, [np]:%v\n", _p, pos, b[pos:endPos], np)
			nodeArr = append(nodeArr, *np)
		} else {
			fmt.Printf("\tisNode=false, _p=%v, pos=%v, payload=% x, [pp]:%v\n", _p, pos, b[pos:endPos], pp)
			primArr = append(primArr, *pp)
		}
	}

	res := &NodePacket{
		NodePackets:      nodeArr,
		PrimitivePackets: primArr,
	}

	return res, endPos, nil
}
