package mark

import (
	"fmt"

	"github.com/yomorun/y3-codec-golang/internal/utils"
)

// Tag represents the Tag of TLV
// MSB used to represent the packet type, 0x80 means a node packet, otherwise is a primitive packet
// Low 7 bits represent Sequence ID, like `key` in JSON format
type Tag struct {
	raw byte
}

// IsNode returns true is MSB is 1.
func (t *Tag) IsNode() bool {
	return t.raw&utils.MSB == utils.MSB
}

// SeqID get the sequence ID, as key in JSON format
func (t *Tag) SeqID() byte {
	//return t.raw & packetutils.DropMSB
	return t.raw & utils.DropMSBArrayFlag
}

// String return string description
func (t *Tag) String() string {
	return fmt.Sprintf("Tag: raw=%4b, SeqID=%v", t.raw, t.SeqID())
}

// NewTag create a NodePacket Tag field
func NewTag(b byte) *Tag {
	return &Tag{raw: b}
}

// IsSlice determine if the current node is a Slice
func (t *Tag) IsSlice() bool {
	return t.raw&utils.SliceFlag == utils.SliceFlag
}

// Raw return the original byte
func (t *Tag) Raw() byte {
	return t.raw
}
