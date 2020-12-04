package codes

import (
	"io"

	y3 "github.com/yomorun/yomo-codec-golang"

	ycodec "github.com/yomorun/yomo-codec-golang/internal/codec"
	"github.com/yomorun/yomo-codec-golang/pkg/spec/encoding"
)

type streamingCodec struct {
	//Packet *y3.NodePacket
	Value []byte

	Length int32
	Size   int32

	Observe byte
	Sbuf    []byte
	Status  decoderStatus

	Matching     [][]byte
	Result       [][]byte
	OriginResult [][]byte

	stick           stickyStatus
	stickyTag       *ycodec.Tag
	stickySize      int32
	stickyLength    int32
	stickyLengthBuf []byte
	stickyTabByte   byte

	proto ProtoCodec
}

func NewStreamingCodec(observe byte) YomoCodec {
	codec := &streamingCodec{
		//Packet: nil,
		Value: make([]byte, 0),

		Observe:      observe,
		Sbuf:         make([]byte, 0),
		Status:       decoderInit,
		Matching:     make([][]byte, 0),
		Result:       make([][]byte, 0),
		OriginResult: make([][]byte, 0),

		stick:           stickyInit,
		stickyTag:       nil,
		stickySize:      0,
		stickyLength:    0,
		stickyLengthBuf: make([]byte, 0),
		stickyTabByte:   byte(0),

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

type stickyStatus uint8

const (
	stickyInit      stickyStatus = 0
	stickyTag       stickyStatus = 1
	stickyLength    stickyStatus = 2
	stickyTagLength stickyStatus = 3
	stickyValue     stickyStatus = 4
)

// Decoder: Collects bytes from buf and decodes them
func (d *streamingCodec) Decoder(buf []byte) {
	key := d.Observe

	var (
		tag       *ycodec.Tag = nil
		size      int32       = 0
		length    int32       = 0
		lengthBuf             = make([]byte, 0)
		curBuf                = make([]byte, 0)
	)

	//fmt.Printf("@110 buf=%v\n", packetutils.FormatBytes(buf))
	for _, c := range buf {
		// tag
		if tag == nil && (d.stick == stickyInit || d.stick == stickyTag || d.stick == stickyTagLength) {
			tag = ycodec.NewTag(c)
			d.Sbuf = append(d.Sbuf, c)
			curBuf = append(curBuf, c)
			if d.Status == decoderInit {
				d.Status = decoderTag
			}
			// sticky handle
			//fmt.Printf("@110 tag.SeqID=%#x\n", tag.SeqID())
			if d.stick == stickyTag || d.stick == stickyTagLength {
				d.stickyTabByte = c
				d.stickyTag = tag
			}
			if d.stick == stickyTagLength {
				d.stick = stickyLength
			}
			continue
		}

		// length
		if size == 0 && (d.stick == stickyInit || d.stick == stickyLength) {
			lengthBuf = append(lengthBuf, c)
			d.Sbuf = append(d.Sbuf, c)
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
			// sticky handle
			//fmt.Printf("@110 size=%v, length=%v\n", size, length)
			if d.stick == stickyLength {
				d.stickyLengthBuf = lengthBuf
				d.stickySize = s
				d.stickyLength = l
				if int32(len(buf))-(1+size) >= 0 {
					d.stick = stickyValue
				}
			}
			continue
		}

		if d.Status == decoderLength {
			d.Status = decoderValue
		}

		//fmt.Printf("@111 stick=%v, buf=%v\n", d.stick, packetutils.FormatBytes(buf))
		//if tag != nil {
		//	fmt.Printf("@111 stick=%v, tag.SeqID=%#x, buf=%v\n", d.stick, tag.SeqID(), packetutils.FormatBytes(buf))
		//}

		if tag != nil && key != tag.SeqID() {
			var newBuf []byte //B:构建新的buf跳到下一个TLV
			if len(buf) > int(1+size+length) {
				//fmt.Printf("#104 tag.SeqID()=%#x, 1+size+length=%v, len(buf)=%v, buf=%v\n",
				//	tag.SeqID(), 1+size+length, len(buf), packetutils.FormatBytes(buf))

				newBuf = buf[1+size+length:]
				d.Sbuf = append(d.Sbuf, buf[(1+size):(1+size+length)]...) //B:把跳过的字节写入Sbuf，不要漏掉
			} else {
				//fmt.Printf("#111 1+size=%v\n", 1+size)
				newBuf = buf[1+size:]
			}

			//fmt.Printf("#101 Tag.SeqID()=%#x, newBuf=%v, Sbuf=%v, curBuf=%v\n",
			//	tag.SeqID(), packetutils.FormatBytes(newBuf), packetutils.FormatBytes(d.Sbuf), packetutils.FormatBytes(curBuf))

			// TODO: 如果产生了粘包情况才需要修改stick值
			// sticky handle
			//fmt.Printf("@333 %v\n", int32(len(newBuf))+int32(len(d.Sbuf)) < (1+d.Size+d.Length))
			if int32(len(newBuf))+int32(len(d.Sbuf)) < (1 + d.Size + d.Length) {
				if len(newBuf) >= 2 {
					d.stick = stickyTagLength //TL,V
				} else if len(newBuf) == 1 {
					//TODO: 未处理切包在T,L与T,LV
					d.stick = stickyTag //T,L
				}
				d.stickyTabByte = byte(0)
				d.stickyLengthBuf = make([]byte, 0)
				d.stickyTag = nil
				d.stickySize = 0
				d.stickyLength = 0
			}

			d.Decoder(newBuf)
			return
		}

		d.Sbuf = append(d.Sbuf, c)
		curBuf = append(curBuf, c)

		//fmt.Printf("#102 Tag.SeqID()=%#x, curBuf=%v, Sbuf=%v, 1+size+length=%v\n",
		//	tag.SeqID(), packetutils.FormatBytes(curBuf), packetutils.FormatBytes(d.Sbuf), 1+size+length)

		// 处理匹配key的值
		if tag != nil && key == tag.SeqID() && int32(len(curBuf)) == 1+size+length {
			//fmt.Printf("#103 Tag.SeqID()=%#x, curBuf=%v, Sbuf=%v\n",
			//	tag.SeqID(), packetutils.FormatBytes(curBuf), packetutils.FormatBytes(d.Sbuf))
			// 当前为匹配上的key，当前的curBuf即为Packet的待解[]byte
			d.Matching = append(d.Matching, curBuf)
			d.Status = decoderMatching
			//debug:临时测试
			//d.reset()
			//return
			//<--
		} else if d.stick == stickyValue && d.stickyTag != nil && key == d.stickyTag.SeqID() && int32(len(curBuf)) == d.stickyLength {
			// `int32(len(curBuf)) == d.stickyLength`只处理了TL,V场景
			// 这里可能有粘包问题
			// 因为是粘包，需要把添加上之前的stickyLengthBuf
			stickyBuf := append(append([]byte{d.stickyTabByte}, d.stickyLengthBuf...), curBuf...)
			d.Matching = append(d.Matching, stickyBuf)
			d.Status = decoderMatching
			//debug:临时测试
			//d.reset()
			//return
			//<--
		}
		//fmt.Printf("@444 stick=%v, stickyTag=%v, curBuf-len=%v, len=%v\n",
		//	d.stick, d.stickyTag, len(curBuf), 1+d.stickySize+d.stickyLength)

		// 收集完根TLV的所有数据后进行分配处理
		if int32(len(d.Sbuf)) == 1+d.Size+d.Length {
			//fmt.Printf("@222 stick=%v, d.Status=%v\n", d.stick, d.Status) //2:decoderLength
			if d.Status == decoderMatching {
				// 如果有matching到被observe的key，Result保存当前根TLV，Matching保存匹配上的TLV，如果匹配上则同时有值
				// Result用于合并数据时使用；Matching用于Read时使用
				d.Result = append(d.Result, d.Sbuf)
				d.OriginResult = append(d.OriginResult, placeholder)
			} else {
				d.OriginResult = append(d.OriginResult, d.Sbuf)
			}

			//重置参数
			d.Status = decoderFinished
			d.reset()
		}
	}
}

// decodeLength: decode length of `V`
func (d *streamingCodec) decodeLength(buf []byte) (length int32, size int32, err error) {
	varCodec := encoding.VarCodec{}
	err = varCodec.DecodePVarInt32(buf, &length)
	size = int32(varCodec.Size)
	return
}

// reset: reset status of the codec
func (d *streamingCodec) reset() {
	d.Length = 0
	d.Size = 0
	d.Sbuf = make([]byte, 0)
	d.Status = decoderInit

	d.stick = stickyInit
	d.stickyTag = nil
	d.stickySize = 0
	d.stickyLength = 0
}

// Read: read and unmarshal data to mold
func (d *streamingCodec) Read(mold interface{}) (interface{}, error) {
	if len(d.Result) == 0 || len(d.Matching) == 0 {
		return nil, nil
	}

	matching := d.Matching[0]
	d.Matching = d.Matching[1:]

	result := d.Result[0]
	d.Result = d.Result[1:]

	//fmt.Printf("#555 %v\n", packetutils.FormatBytes(matching))

	proto := NewProtoCodec(d.Observe)
	if proto.IsStruct(mold) {
		//err := proto.UnmarshalStruct(result, mold)
		//if err != nil {
		//	return nil, err
		//}
		// #2
		//packet, _, err := y3.DecodeNodePacket(matching)
		//if err != nil {
		//	return nil, err
		//}
		//err = packetstructure.Decode(packet, mold)
		//if err != nil {
		//	return nil, err
		//}

		err := proto.UnmarshalStructNative(matching, mold)
		if err != nil {
			return nil, err
		}

	} else {
		//err := proto.UnmarshalBasic(matching, &mold)
		//if err != nil {
		//	return nil, err
		//}
		err := proto.UnmarshalBasicNative(matching, &mold)
		if err != nil {
			return nil, err
		}
	}

	// for Encoder::merge
	d.Value = result
	//d.Packet = matching

	return mold, nil
}

// Write: write interface to stream
func (d *streamingCodec) Write(w io.Writer, T interface{}, mold interface{}) (int, error) {
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
func (d *streamingCodec) Encoder(buf []byte, mold interface{}) ([]byte, error) {
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
func (d *streamingCodec) isDecoder(buf []byte) bool {
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
func (d *streamingCodec) mergePacket(source *y3.NodePacket, key byte, value []byte, mold interface{}) ([]byte, error) {
	np := d.copyPacket(nil, source, key, value)
	buf := np.Encode()
	return buf, nil
}

// copyPacket: copy packet
func (d *streamingCodec) copyPacket(root *y3.NodePacketEncoder, source *y3.NodePacket, key byte, value []byte) *y3.NodePacketEncoder {
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
func (d *streamingCodec) Refresh(w io.Writer) (int, error) {
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
