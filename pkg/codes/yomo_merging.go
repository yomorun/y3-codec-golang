package codes

import (
	"io"

	y3 "github.com/yomorun/yomo-codec-golang"

	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"
)

type mergingCodec struct {
	Value []byte

	Length int32
	Size   int32

	Observe   byte
	SourceBuf []byte
	Status    decoderStatus

	Matching     [][]byte
	Result       [][]byte
	OriginResult [][]byte

	CollectedStatus    collectedStatus
	CollectedTag       *ycodec.Tag
	CollectedSize      int32
	CollectedLength    int32
	CollectedLengthBuf []byte
	CollectedBuffer    []byte
	CollectedResult    [][]byte

	proto ProtoCodec
}

func NewMergingCodec(observe byte) YomoCodec {
	codec := &mergingCodec{
		Value: make([]byte, 0),

		Observe:      observe,
		SourceBuf:    make([]byte, 0),
		Status:       decoderInit,
		Matching:     make([][]byte, 0),
		Result:       make([][]byte, 0),
		OriginResult: make([][]byte, 0),

		CollectedStatus:    collectedInit,
		CollectedTag:       nil,
		CollectedSize:      0,
		CollectedLength:    0,
		CollectedLengthBuf: make([]byte, 0),
		CollectedBuffer:    make([]byte, 0),
		CollectedResult:    make([][]byte, 0),

		proto: NewProtoCodec(observe),
	}
	return codec
}

// Decoder: Collects bytes from buf and decodes them
func (d *mergingCodec) Decoder(buf []byte) {

	for _, c := range buf {
		// tag
		if d.CollectedStatus == collectedInit && d.CollectedTag == nil {
			d.CollectedTag = ycodec.NewTag(c)
			d.CollectedStatus = collectedTag
			continue
		}

		// length
		if d.CollectedStatus == collectedTag || d.CollectedStatus == collectedLength {
			d.CollectedLengthBuf = append(d.CollectedLengthBuf, c)
			l, s, err := d.decodeLength(d.CollectedLengthBuf)
			if err != nil {
				d.CollectedStatus = collectedLength
				continue
			}
			d.CollectedSize = s
			d.CollectedLength = l
			d.CollectedStatus = collectedBody
			continue
		}

		if d.CollectedStatus == collectedBody {
			d.CollectedBuffer = append(d.CollectedBuffer, d.CollectedTag.Raw())
			d.CollectedBuffer = append(d.CollectedBuffer, d.CollectedLengthBuf...)
			d.CollectedStatus = collectedCaching
		}

		if d.CollectedStatus == collectedCaching {
			if !d.isCollectCompleted() {
				d.CollectedBuffer = append(d.CollectedBuffer, c)
			}
			if d.isCollectCompleted() {
				d.CollectedResult = append(d.CollectedResult, d.CollectedBuffer)
				d.CollectedStatus = collectedFinished
				d.resetCollector()
			}
		}

	}

	if len(d.CollectedResult) > 0 {
		for i := 0; i < len(d.CollectedResult)+1; i++ {
			result := d.CollectedResult[0]
			d.CollectedResult = d.CollectedResult[1:]
			d.decode(result)
		}
	}
}

func (d *mergingCodec) isCollectCompleted() bool {
	return (1+d.CollectedSize+d.CollectedLength)-int32(len(d.CollectedBuffer)) == 0
}

func (d *mergingCodec) resetCollector() {
	d.CollectedStatus = collectedInit
	d.CollectedTag = nil
	d.CollectedSize = 0
	d.CollectedLength = 0
	d.CollectedLengthBuf = make([]byte, 0)
	d.CollectedBuffer = make([]byte, 0)
}

func (d *mergingCodec) decode(buf []byte) {
	key := d.Observe

	var (
		tag       *ycodec.Tag = nil
		size      int32       = 0
		length    int32       = 0
		lengthBuf             = make([]byte, 0)
		curBuf                = make([]byte, 0)
	)

	for _, c := range buf {
		// tag
		if tag == nil {
			tag = ycodec.NewTag(c)
			d.SourceBuf = append(d.SourceBuf, c)
			curBuf = append(curBuf, c)
			if d.Status == decoderInit {
				d.Status = decoderTag
			}
			continue
		}

		// length
		if size == 0 {
			lengthBuf = append(lengthBuf, c)
			d.SourceBuf = append(d.SourceBuf, c)
			curBuf = append(curBuf, c)
			l, s, err := d.decodeLength(lengthBuf)
			if err != nil {
				continue
			}
			size = s
			length = l
			if d.Status == decoderTag {
				d.Size = s
				d.Length = l
				d.Status = decoderLength
			}
			continue
		}

		if tag != nil && key != tag.SeqID() {
			var newBuf []byte
			if len(buf) > int(1+size+length) {
				newBuf = buf[1+size+length:]
				d.SourceBuf = append(d.SourceBuf, buf[(1+size):(1+size+length)]...)
			} else {
				newBuf = buf[1+size:]
			}

			d.decode(newBuf)
			return
		}

		if d.Status == decoderLength {
			d.Status = decoderValue
		}

		d.SourceBuf = append(d.SourceBuf, c)
		curBuf = append(curBuf, c)

		if tag != nil && key == tag.SeqID() && int32(len(curBuf)) == 1+size+length {
			d.Matching = append(d.Matching, curBuf)
			d.Status = decoderMatching
		}

		if int32(len(d.SourceBuf)) == 1+d.Size+d.Length {
			if d.Status == decoderMatching {
				d.Result = append(d.Result, d.SourceBuf)
				d.OriginResult = append(d.OriginResult, placeholder)
			} else {
				d.OriginResult = append(d.OriginResult, d.SourceBuf)
			}

			d.Status = decoderFinished
			d.reset()
		}
	}
}

// decodeLength: decode length of `V`
func (d *mergingCodec) decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

// reset: reset status of the codec
func (d *mergingCodec) reset() {
	d.Length = 0
	d.Size = 0
	d.SourceBuf = make([]byte, 0)
	d.Status = decoderInit

}

// Read: read and unmarshal data to mold
func (d *mergingCodec) Read(mold interface{}) (interface{}, error) {
	if len(d.Result) == 0 || len(d.Matching) == 0 {
		return nil, nil
	}

	matching := d.Matching[0]
	d.Matching = d.Matching[1:]

	result := d.Result[0]
	d.Result = d.Result[1:]

	info := &MoldInfo{Mold: mold}
	err := d.proto.Unmarshal(matching, info)
	if err != nil {
		return nil, err
	}
	mold = info.Mold

	// for Encoder::merge
	d.Value = result

	return mold, nil
}

// Write: write interface to stream
func (d *mergingCodec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
	buf, err := d.proto.MarshalNative(T)
	if err != nil {
		logger.Errorf("Write::MarshalNative error:%v", err)
		return 0, err
	}

	data, err := d.Encoder(buf, mold)
	if err != nil {
		return 0, err
	}

	return w.Write(data)
}

// Encoder: encode []byte of T, and merge them to the original result
func (d *mergingCodec) Encoder(buf []byte, mold interface{}) ([]byte, error) {
	result := make([]byte, 0)
	index := 0
	for _, data := range d.OriginResult {
		index = index + 1
		if d.isDecoder(data) {
			source, _, _ := y3.DecodeNodePacket(d.Value)

			key := d.Observe
			_buf, err := d.mergePacket(source, key, buf, mold)
			if err != nil {
				return nil, err
			}

			d.Value = make([]byte, 0)
			result = append(result, _buf...)
			break
		} else {
			result = append(result, data...)
		}
	}

	d.OriginResult = d.OriginResult[index:]

	return result, nil
}

// isDecoder: is placeholder?
func (d *mergingCodec) isDecoder(buf []byte) bool {
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

// mergePacket: merge packet
func (d *mergingCodec) mergePacket(source *y3.NodePacket, key byte, value []byte, mold interface{}) ([]byte, error) {
	np := d.copyPacket(nil, source, key, value)
	buf := np.Encode()
	return buf, nil
}

// copyPacket: copy packet
func (d *mergingCodec) copyPacket(root *y3.NodePacketEncoder, source *y3.NodePacket, key byte, value []byte) *y3.NodePacketEncoder {
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
			np := d.copyPacket(temp, &n, key, value)
			root.AddNodePacket(np)
		}
	}

	return root
}

// Refresh: refresh the OriginResult to stream
func (d *mergingCodec) Refresh(w io.Writer) (int, error) {
	if len(d.OriginResult) == 0 {
		return 0, nil
	}
	originResult := d.OriginResult[0]
	if !d.isDecoder(originResult) {
		d.OriginResult = d.OriginResult[1:]
		return w.Write(originResult)
	}
	return 0, nil
}
