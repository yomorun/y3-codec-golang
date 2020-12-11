package codes

import (
	"io"

	y3 "github.com/yomorun/yomo-codec-golang"
	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/pkg/packetutils"
	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"
)

// collectingCodec: Implementation of the YomoCodec Interface
type collectingCodec struct {
	Value *y3.NodePacket

	Tag *ycodec.Tag

	LengthBuf []byte
	Length    int32
	Size      int32

	Observe byte
	Sbuf    []byte

	Result       []*y3.NodePacket
	OriginResult [][]byte

	proto ProtoCodec
}

func NewCollectingCodec(observe string) YomoCodec {
	ob := packetutils.KeyOf(observe)
	codec := &collectingCodec{
		Value:        nil,
		Tag:          nil,
		LengthBuf:    make([]byte, 0),
		Length:       0,
		Size:         0,
		Sbuf:         make([]byte, 0),
		Observe:      ob,
		Result:       make([]*y3.NodePacket, 0),
		OriginResult: make([][]byte, 0),
		proto:        NewProtoCodec(ob),
	}

	return codec
}

func (codec *collectingCodec) Decoder(buf []byte) {
	key := codec.Observe
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

		if int32(len(codec.Sbuf)) == 1+codec.Size+codec.Length {
			packet, _, err := y3.DecodeNodePacket(codec.Sbuf)
			if err != nil {
				logger.Errorf("::Decoder DecodeNodePacket error:%v", err)
				codec.reset()
				continue
			}

			var result *y3.NodePacket
			originResult := codec.Sbuf
			codec.reset()

			//matching
			var matching = false
			flag, _, _ := packetutils.MatchingKey(key, packet)
			if flag {
				matching = true
				result = packet
			}

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
func (codec *collectingCodec) reset() {
	codec.Tag = nil
	codec.LengthBuf = make([]byte, 0)
	codec.Length = 0
	codec.Size = 0
	codec.Sbuf = make([]byte, 0)
}

// decodeLength: decode length of `V`
func (codec *collectingCodec) decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

// Read: read and unmarshal data to mold
func (codec *collectingCodec) Read(mold interface{}) (interface{}, error) {
	if len(codec.Result) == 0 {
		return nil, nil
	}

	result := codec.Result[0]
	codec.Result = codec.Result[1:]

	proto := codec.proto
	if proto.IsStruct(mold) {
		err := proto.UnmarshalStructByNodePacket(result, mold)
		if err != nil {
			return nil, err
		}
	} else {
		err := proto.UnmarshalBasicByNodePacket(result, &mold)
		if err != nil {
			return nil, err
		}
	}

	// for Encoder::merge
	codec.Value = result

	return mold, nil
}

// Write: write interface to stream
func (codec *collectingCodec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
	proto := codec.proto
	buf, err := proto.MarshalNative(T)
	if err != nil {
		logger.Errorf("Write::MarshalNative error:%v", err)
		return 0, err
	}

	data, err := codec.Encoder(buf, mold)
	if err != nil {
		return 0, err
	}

	return w.Write(data)
}

// Encoder: encode []byte of T, and merge them to the original result
func (codec *collectingCodec) Encoder(buf []byte, mold interface{}) ([]byte, error) {
	result := make([]byte, 0)
	index := 0
	for _, data := range codec.OriginResult {
		index = index + 1
		if codec.isDecoder(data) {
			source := codec.Value

			key := codec.Observe
			_buf, err := codec.mergePacket(source, key, buf, mold)
			if err != nil {
				return nil, err
			}

			codec.Value = nil
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
func (codec *collectingCodec) mergePacket(source *y3.NodePacket, key byte, value []byte, mold interface{}) ([]byte, error) {
	np := codec.copyPacket(nil, source, key, value)
	buf := np.Encode()
	return buf, nil
}

// copyPacket: copy packet
func (codec *collectingCodec) copyPacket(root *y3.NodePacketEncoder, source *y3.NodePacket, key byte, value []byte) *y3.NodePacketEncoder {
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
func (codec *collectingCodec) Refresh(w io.Writer) (int, error) {
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
func (codec *collectingCodec) isDecoder(buf []byte) bool {
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
