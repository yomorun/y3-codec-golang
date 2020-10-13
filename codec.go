package y3

import (
	"encoding/hex"
	"io"
	"strings"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
	encoding "github.com/yomorun/yomo-codec-golang/pkg"
)

var (
	placeholder = []byte{0, 1, 2, 3}
	logger      = utils.Logger.WithPrefix(utils.DefaultLogger, "Codec")
)

type Codec struct {
	Value []byte

	Tag *ycodec.Tag

	LengthBuf []byte
	Length    int32
	Size      int32

	Observe string
	Sbuf    []byte

	Result       [][]byte
	OriginResult [][]byte
}

func NewCodec(observe string) *Codec {
	codec := &Codec{
		Value:        make([]byte, 0),
		Tag:          nil,
		LengthBuf:    make([]byte, 0),
		Length:       0,
		Size:         0,
		Sbuf:         make([]byte, 0),
		Observe:      observe,
		Result:       make([][]byte, 0),
		OriginResult: make([][]byte, 0),
	}

	return codec
}

func (codec *Codec) Decoder(buf []byte) {
	key := keyOf(codec.Observe)
	for _, c := range buf {
		// tag
		if codec.Tag == nil {
			codec.Tag = ycodec.NewTag(c)
			codec.Sbuf = append(codec.Sbuf, c)
			continue
		}

		// length
		if codec.Size == 0 {
			codec.LengthBuf = append(codec.LengthBuf, c)
			codec.Sbuf = append(codec.Sbuf, c)
			length, size, err := decodeLength(codec.LengthBuf)
			if err != nil {
				continue
			}
			codec.Length = length
			codec.Size = size
			continue
		}

		codec.Sbuf = append(codec.Sbuf, c)

		// buf end, then handle Sbuf
		if int32(len(codec.Sbuf)) == 1+codec.Size+codec.Length {
			// Decode Packet from Sbuf
			packet, _, err := DecodeNodePacket(codec.Sbuf)
			if err != nil {
				logger.Errorf("::Decoder DecodeNodePacket error:%v", err)
				codec.reset()
				continue
			}

			// temp save Sbuf and reset
			result := make([]byte, 0)
			originResult := codec.Sbuf
			codec.reset()

			//matching
			var matching = false
			flag, _, _ := matchingKey(key, packet)
			if flag || []byte("*")[0] == key {
				matching = true
				result = originResult
			}

			// save to result
			if matching {
				codec.Result = append(codec.Result, result)
				codec.OriginResult = append(codec.OriginResult, placeholder)
			} else {
				codec.OriginResult = append(codec.OriginResult, originResult)
			}

		}
	}
}

func (codec *Codec) reset() {
	codec.Tag = nil
	codec.LengthBuf = make([]byte, 0)
	codec.Length = 0
	codec.Size = 0
	codec.Sbuf = make([]byte, 0)
}

func matchingKey(key byte, node *NodePacket) (flag bool, isNode bool, packet interface{}) {
	if len(node.PrimitivePackets) > 0 {
		for _, p := range node.PrimitivePackets {
			if key == p.tag.SeqID() {
				return true, false, p
			}
		}
	}

	if len(node.NodePackets) > 0 {
		for _, n := range node.NodePackets {
			if key == n.tag.SeqID() {
				return true, true, n
			}
			//return matchingKey(key, &n)
			flag, isNode, packet = matchingKey(key, &n)
			if flag {
				return
			}
		}
	}

	return false, false, nil
}

func keyOf(hexStr string) byte {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = strings.TrimPrefix(hexStr, "0x")
	} else if strings.HasPrefix(hexStr, "0X") {
		hexStr = strings.TrimPrefix(hexStr, "0X")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		logger.Errorf("hex.DecodeString error: %v", err)
		return 0xff
	}

	if len(data) == 0 {
		logger.Errorf("hex.DecodeString data is []")
		return 0xff
	}

	return data[0]
}

func decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

func (codec *Codec) Read(mold interface{}) (interface{}, error) {
	if len(codec.Result) == 0 {
		return nil, nil
	}

	result := codec.Result[0]
	codec.Result = codec.Result[1:]

	// take value from node
	err := codec.Unmarshal(result, &mold)
	if err != nil {
		return nil, err
	}

	// for Encoder::merge
	codec.Value = result

	return mold, nil
}

func (codec *Codec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
	// #1. mold --> NodePacket
	// #2. merge NodePacket --> codec.Value NodePacket
	// #3. NodePacket --> []byte
	// #4. Write []byte

	buf, err := codec.Marshal(T)
	if err != nil {
		logger.Errorf("Write::Marshal error:%v", err)
		return 0, err
	}

	data, err := codec.Encoder(buf, mold)
	if err != nil {
		return 0, err
	}

	return w.Write(data)
}

func (codec *Codec) Encoder(buf []byte, mold interface{}) ([]byte, error) {
	result := make([]byte, 0)
	index := 0
	for _, data := range codec.OriginResult {
		index = index + 1
		if isDecoder(data) {
			source, _, _ := DecodeNodePacket(codec.Value)

			key := keyOf(codec.Observe)
			_buf, err := mergePacket(source, key, buf, mold)
			if err != nil {
				return nil, err
			}

			codec.Value = make([]byte, 0)
			result = append(result, _buf...)
			break
		} else {
			result = append(result, data...)
		}
	}

	codec.OriginResult = codec.OriginResult[index:]

	return result, nil
}

func mergePacket(source *NodePacket, key byte, value []byte, mold interface{}) ([]byte, error) {
	np := copyPacket(nil, source, key, value)
	buf := np.Encode()
	return buf, nil
}

func copyPacket(root *NodePacketEncoder, source *NodePacket, key byte, value []byte) *NodePacketEncoder {
	if root == nil {
		root = NewNodePacketEncoder(int(source.tag.SeqID()))
	}

	if len(source.PrimitivePackets) > 0 {
		for _, p := range source.PrimitivePackets {
			temp := NewPrimitivePacketEncoder(int(p.tag.SeqID()))
			if p.tag.SeqID() == key {
				temp.SetBytes(value)
			} else {
				temp.SetBytes(p.ToBytes())
			}

			root.AddPrimitivePacket(temp)
		}
	}

	if len(source.NodePackets) > 0 {
		for _, n := range source.NodePackets {
			var temp *NodePacketEncoder
			if n.IsArray() {
				temp = NewNodeArrayPacketEncoder(int(n.SeqID()))
			} else {
				temp = NewNodePacketEncoder(int(n.SeqID()))
			}

			if n.tag.SeqID() == key {
				// replace node
				temp.AddBytes(value)
				root.AddNodePacket(temp)
				continue
			}
			np := copyPacket(temp, &n, key, value)
			root.AddNodePacket(np)
		}
	}

	return root
}

func (codec *Codec) Refresh(w io.Writer) (int, error) {
	if len(codec.OriginResult) == 0 {
		return 0, nil
	}
	originResult := codec.OriginResult[0]
	if !isDecoder(originResult) {
		codec.OriginResult = codec.OriginResult[1:]
		return w.Write(originResult)
	}
	return 0, nil
}

func isDecoder(buf []byte) bool {
	if len(buf) != len(placeholder) {
		return false
	}
	for i, v := range placeholder {
		if buf[i] != v {
			return false
		}
	}
	return true
}
