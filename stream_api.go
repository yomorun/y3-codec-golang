package y3

import (
	"io"
	"log"

	"github.com/yomorun/y3-codec-golang/pkg/encoding"
	"github.com/yomorun/y3-codec-golang/pkg/spec"
)

// StreamDecoder decode Y3 Packet from a io.Reader
type StreamDecoder struct {
	errState bool
	tagbuf   []byte
	lenbuf   []byte
	valbuf   []byte
	r        io.Reader
	state    string
	len      int
	callback func(*spec.Packet)
}

// NewStreamDecoder return a stream decoder
func NewStreamDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{
		errState: false,
		r:        r,
		state:    "Nil",
	}
}

// OnPacket trigger callback once Y3 packet parsed out
func (sd *StreamDecoder) OnPacket(f func(*spec.Packet)) {
	sd.callback = f
}

// Start the parser
func (sd *StreamDecoder) Start() {
	// buffer
	tmp := make([]byte, 1)
	for {
		n, err := sd.r.Read(tmp)
		if err != nil {
			log.Printf("io err: tmp=[% X]", tmp)
			sd.reset(err)
			break
		}
		log.Printf("Recieved: n=%d, tmp=[% X]", n, tmp[:n])
		for _, v := range tmp[:n] {
			sd.fill(v)
		}
	}
}

func (sd *StreamDecoder) fill(b byte) error {
	log.Printf("-> fill b=[% X], state=%s", b, sd.state)
	switch sd.state {
	case "Nil":
		sd.state = "TS"
		sd.fill(b)
	case "TS":
		sd.tagbuf = append(sd.tagbuf, b)
		if b&0x81 != 0x81 {
			// over of tag
			sd.state = "LS"
			log.Printf("Parsed Out Tag, tagbuf=[% X]", sd.tagbuf)
			return nil
		}
	case "LS":
		sd.lenbuf = append(sd.lenbuf, b)
		if b&0x81 != 0x81 {
			// over of len, start parse as PVarUInt64 value
			var len uint64
			codec := encoding.VarCodec{}
			err := codec.DecodePVarUInt64(sd.lenbuf, &len)
			if err != nil {
				sd.errState = true
				panic(err)
			} else {
				sd.len = int(len)
				log.Printf("Parsed Out len=%d, lenbuf=[% X]", sd.len, sd.lenbuf)
			}
			if sd.len == 0 {
				// reset state if zero-len packet
				log.Printf("[%s] Parsed Out valbuf=EMPTY", sd.state)
				// make a Packet object
				sd.fullfiled()
				sd.state = "Nil"
			}
			// update state
			sd.state = "VS"
			return nil
		}
	case "VS":
		sd.valbuf = append(sd.valbuf, b)
		if len(sd.valbuf) == sd.len {
			log.Printf("[%s] Parsed Out valbuf=[% X]", sd.state, sd.valbuf)
			// make a Packet object
			sd.fullfiled()
			// reset state
			sd.state = "Nil"
		}
	}
	return nil
}

func (sd *StreamDecoder) fullfiled() {
	buf := append(sd.tagbuf, sd.lenbuf...)
	buf = append(buf, sd.valbuf...)

	p, err := spec.FromBytes(buf)
	if err != nil {
		panic(err)
	}
	log.Printf("--> Fullfiled p=%v", p)
	sd.callback(p)

	sd.reset(nil)
}

func (sd *StreamDecoder) reset(err error) {
	if err != nil {
		log.Printf("[RESET] cause of error: %s", err.Error())
	}
	sd.errState = false
	sd.tagbuf = make([]byte, 0)
	sd.lenbuf = make([]byte, 0)
	sd.valbuf = make([]byte, 0)
	sd.state = "Nil"
}
