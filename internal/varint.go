package varint

import (
	"encoding/hex"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// Varint 定义了一种描述整数的方法, 它是变长类型
//
// 特点：
// 1. 如果当前byte的highest bit是1，则下一个byte也是Varint的一部分。而这个highest bit我们称为MSB
// 2. 如果当前byte的highest bit是0，则该byte是Varint的最后一个部分
// 3. 去掉每一个byte的最高位，按小端序计算
// 4. 使用zigzag算法计算出结果
type Varint struct {
	raw []byte
}

// Decoder 用于解码
type Decoder struct {
	raw    []byte
	logger utils.Logger
}

// NewDecoder return a decoder for parsing `Varint`
func NewDecoder(buf []byte) (dec *Decoder, len int) {
	// 因为是可变类型，所以此时还不知道有buf中有多少bytes是Varint使用的
	// 根据Varint类型的特点来寻找Varint的结束位置
	// 并将属于Varint类型的buffer传入Decoder中
	len = 0
	raw := make([]byte, 0)
	for _, b := range buf {
		len++
		raw = append(raw, b)
		if b&0x80 != 0x80 {
			break
		}
	}

	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "varint::Decode")
	return &Decoder{raw, logger}, len
}

// Decode parse bytes and returns the `uint64` value
func (d *Decoder) Decode() (int64, error) {
	d.logger.Debugf("raw=%s", hex.Dump(d.raw))
	var val uint64
	for i, v := range d.raw {
		val |= (uint64(v & 0x7F)) << (i * 7)
	}
	res := zigzagDecode(val)
	return res, nil
}

/// zigzag from : https://developers.google.com/protocol-buffers/docs/encoding#signed-integers
func zigzagEncode(from int64) uint64 {
	return uint64((from << 1) ^ (from >> 63))
}

/// zigzag from : https://developers.google.com/protocol-buffers/docs/encoding#signed-integers
func zigzagDecode(from uint64) int64 {
	return int64((from >> 1) ^ uint64(-(int64(from & 1))))
}

// ParseVarintLength parse length as a varint type
func ParseVarintLength(b []byte, startPos int) (val int64, lengthOfLenBuffer int, err error) {
	dec, lengthOfLenBuffer := NewDecoder(b[startPos:])
	val, err = dec.Decode()
	if err != nil {
		return 0, 0, err
	}
	return val, lengthOfLenBuffer, nil
}
