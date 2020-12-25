package codes

import (
	"io"
	"sync"

	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"

	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
)

type streamingCodec struct {
	inform       chan bool
	enableInform bool

	Length int32
	Size   int32

	Observe   byte
	SourceBuf []byte
	Status    decoderStatus

	Matching          [][]byte
	matchingMutex     sync.Mutex
	OriginResult      [][]byte
	OriginResultMutex sync.Mutex

	CollectedStatus    collectedStatus
	CollectedTag       *ycodec.Tag
	CollectedSize      int32
	CollectedLength    int32
	CollectedLengthBuf []byte
	CollectedBuffer    []byte
	CollectedResult    [][]byte

	proto ProtoCodec
}

func NewStreamingCodec(observe byte) (YomoCodec, <-chan bool) {
	codec := &streamingCodec{
		inform:       make(chan bool, 10),
		enableInform: true,

		Observe:      observe,
		SourceBuf:    make([]byte, 0),
		Status:       decoderInit,
		Matching:     make([][]byte, 0),
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
	return codec, codec.inform
}

func NewStreamingCodecNoInform(observe byte) YomoCodec {
	codec := &streamingCodec{
		inform:       make(chan bool, 10),
		enableInform: false,

		Observe:      observe,
		SourceBuf:    make([]byte, 0),
		Status:       decoderInit,
		Matching:     make([][]byte, 0),
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

type decoderStatus uint8

const (
	decoderInit     decoderStatus = 0
	decoderTag      decoderStatus = 1
	decoderLength   decoderStatus = 2
	decoderValue    decoderStatus = 3
	decoderMatching decoderStatus = 4
	decoderFinished decoderStatus = 5
)

type collectedStatus uint8

const (
	collectedInit     collectedStatus = 0
	collectedTag      collectedStatus = 1
	collectedLength   collectedStatus = 2
	collectedBody     collectedStatus = 3
	collectedCaching  collectedStatus = 4
	collectedFinished collectedStatus = 5
)

// Decoder: Collects bytes from buf and decodes them
func (d *streamingCodec) Decoder(buf []byte) {
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

func (d *streamingCodec) isCollectCompleted() bool {
	return (1+d.CollectedSize+d.CollectedLength)-int32(len(d.CollectedBuffer)) == 0
}

func (d *streamingCodec) resetCollector() {
	d.CollectedStatus = collectedInit
	d.CollectedTag = nil
	d.CollectedSize = 0
	d.CollectedLength = 0
	d.CollectedLengthBuf = make([]byte, 0)
	d.CollectedBuffer = make([]byte, 0)
}

func (d *streamingCodec) decode(buf []byte) {
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
			d.matchingMutex.Lock()
			d.Matching = append(d.Matching, curBuf)
			d.matchingMutex.Unlock()
			if d.enableInform {
				d.inform <- true
			}
			d.Status = decoderMatching
		}

		if int32(len(d.SourceBuf)) == 1+d.Size+d.Length {
			if d.Status == decoderMatching {
				//d.Result = append(d.Result, d.SourceBuf)
				d.OriginResultMutex.Lock()
				d.OriginResult = append(d.OriginResult, placeholder)
				d.OriginResultMutex.Unlock()
			} else {
				d.OriginResultMutex.Lock()
				d.OriginResult = append(d.OriginResult, d.SourceBuf)
				d.OriginResultMutex.Unlock()
			}

			d.Status = decoderFinished
			d.reset()
		}
	}
}

func (d *streamingCodec) decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

func (d *streamingCodec) reset() {
	d.Length = 0
	d.Size = 0
	d.SourceBuf = make([]byte, 0)
	d.Status = decoderInit
}

// Read: read and unmarshal data to mold
func (d *streamingCodec) Read(mold interface{}) (interface{}, error) {
	if len(d.Matching) == 0 {
		return nil, nil
	}
	d.matchingMutex.Lock()
	matching := d.Matching[0]
	d.Matching = d.Matching[1:]
	d.matchingMutex.Unlock()

	info := &MoldInfo{Mold: mold}
	err := d.proto.Unmarshal(matching, info)
	if err != nil {
		return nil, err
	}
	mold = info.Mold

	return mold, nil
}

// Write: write interface to stream
func (d *streamingCodec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
	proto := NewProtoCodec(d.Observe)
	result, _ := proto.Marshal(T)
	return w.Write(result)
}

func (d *streamingCodec) isPlaceholder(buf []byte) bool {
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

// Refresh: refresh the OriginResult to stream
func (d *streamingCodec) Refresh(w io.Writer) (int, error) {
	if len(d.OriginResult) == 0 {
		return 0, nil
	}

	d.OriginResultMutex.Lock()
	originResult := d.OriginResult[0]
	d.OriginResult = d.OriginResult[1:]
	d.OriginResultMutex.Unlock()

	if !d.isPlaceholder(originResult) {
		return w.Write(originResult)
	}
	return 0, nil
}
