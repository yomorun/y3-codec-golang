package codes

import (
	"io"

	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"

	y3 "github.com/yomorun/yomo-codec-golang"

	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"

	"github.com/yomorun/yomo-codec-golang/internal/utils"

	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
)

var (
	placeholder = []byte{0, 1, 2, 3}
	logger      = utils.Logger.WithPrefix(utils.DefaultLogger, "yomoCodec")
)

// YomoCodec: codec interface for yomo
type YomoCodec interface {
	Decoder(buf []byte)
	Read(mold interface{}) (interface{}, error) // TODO: 考虑掉返回interface{}
	Write(w io.Writer, T interface{}, mold interface{}) (int, error)
	Refresh(w io.Writer) (int, error)
}

// yomoCodec: Implementation of the YomoCodec Interface
type yomoCodec struct {
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

func NewCodec(observe string) YomoCodec {
	codec := &yomoCodec{
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

// Decoder: Collects bytes from buf and decodes them
func (codec *yomoCodec) Decoder(buf []byte) {
	key := packetutils.KeyOf(codec.Observe)
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
			length, size, err := codec.decodeLength(codec.LengthBuf)
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
			packet, _, err := y3.DecodeNodePacket(codec.Sbuf)
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
			flag, _, _ := packetutils.MatchingKey(key, packet)
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

// reset: reset status of the codec
func (codec *yomoCodec) reset() {
	codec.Tag = nil
	codec.LengthBuf = make([]byte, 0)
	codec.Length = 0
	codec.Size = 0
	codec.Sbuf = make([]byte, 0)
}

// decodeLength: decode length of `V`
func (codec *yomoCodec) decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

// Read: read and unmarshal data to mold
func (codec *yomoCodec) Read(mold interface{}) (interface{}, error) {
	if len(codec.Result) == 0 {
		return nil, nil
	}

	result := codec.Result[0]
	codec.Result = codec.Result[1:]

	proto := NewProtoCodec(codec.Observe)
	if proto.IsStruct(mold) {
		err := proto.UnmarshalStruct(result, mold)
		if err != nil {
			return nil, err
		}
	} else {
		err := proto.UnmarshalBasic(result, &mold)
		if err != nil {
			return nil, err
		}
	}

	// for Encoder::merge
	codec.Value = result

	return mold, nil
}

// Write: write interface to stream
func (codec *yomoCodec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
	// #1. mold --> NodePacket
	// #2. merge NodePacket --> codec.Value NodePacket
	// #3. NodePacket --> []byte
	// #4. Write []byte
	proto := NewProtoCodec(codec.Observe)
	buf, err := proto.MarshalNoWrapper(T)
	if err != nil {
		logger.Errorf("Write::MarshalNoWrapper error:%v", err)
		return 0, err
	}

	data, err := codec.Encoder(buf, mold)
	if err != nil {
		return 0, err
	}

	return w.Write(data)
}

// Encoder: encode []byte of T, and merge them to the original result
func (codec *yomoCodec) Encoder(buf []byte, mold interface{}) ([]byte, error) {
	result := make([]byte, 0)
	index := 0
	for _, data := range codec.OriginResult {
		index = index + 1
		if codec.isDecoder(data) {
			source, _, _ := y3.DecodeNodePacket(codec.Value)

			key := packetutils.KeyOf(codec.Observe)
			_buf, err := codec.mergePacket(source, key, buf, mold)
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

// mergePacket: merge packet
func (codec *yomoCodec) mergePacket(source *y3.NodePacket, key byte, value []byte, mold interface{}) ([]byte, error) {
	np := codec.copyPacket(nil, source, key, value)
	buf := np.Encode()
	return buf, nil
}

// copyPacket: copy packet
func (codec *yomoCodec) copyPacket(root *y3.NodePacketEncoder, source *y3.NodePacket, key byte, value []byte) *y3.NodePacketEncoder {
	if root == nil {
		root = y3.NewNodePacketEncoder(int(source.SeqID()))
	}

	if len(source.PrimitivePackets) > 0 {
		for _, p := range source.PrimitivePackets {
			temp := y3.NewPrimitivePacketEncoder(int(p.SeqID()))
			if p.SeqID() == key {
				temp.SetBytes(value)
			} else {
				temp.SetBytes(p.ToBytes())
			}

			root.AddPrimitivePacket(temp)
		}
	}

	if len(source.NodePackets) > 0 {
		for _, n := range source.NodePackets {
			var temp *y3.NodePacketEncoder
			if n.IsArray() {
				temp = y3.NewNodeArrayPacketEncoder(int(n.SeqID()))
			} else {
				temp = y3.NewNodePacketEncoder(int(n.SeqID()))
			}

			if n.SeqID() == key {
				// replace node
				temp.AddBytes(value)
				root.AddNodePacket(temp)
				continue
			}
			np := codec.copyPacket(temp, &n, key, value)
			root.AddNodePacket(np)
		}
	}

	return root
}

// Refresh: refresh the OriginResult to stream
func (codec *yomoCodec) Refresh(w io.Writer) (int, error) {
	if len(codec.OriginResult) == 0 {
		return 0, nil
	}
	originResult := codec.OriginResult[0]
	if !codec.isDecoder(originResult) {
		codec.OriginResult = codec.OriginResult[1:]
		return w.Write(originResult)
	}
	return 0, nil
}

// isDecoder: is placeholder?
func (codec *yomoCodec) isDecoder(buf []byte) bool {
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
