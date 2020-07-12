package codec

import (
	"errors"
	"fmt"

	"github.com/yomorun/yomo-codec-golang/internal/utils"
)

// NodeTag represents the Tag of TLTV
type NodeTag struct {
	raw byte
}

// SeqID 获取Key的顺序ID
func (t *NodeTag) SeqID() byte {
	return t.raw & utils.DropMSB
}

func (t *NodeTag) String() string {
	return fmt.Sprintf("Tag: raw=%4b, SeqID=%v", t.raw, t.SeqID())
}

// NewNodeTag create a NodePacket Tag field
func NewNodeTag(b byte) (p *NodeTag, err error) {
	// 最高位始终为1
	if b&utils.MSB != utils.MSB {
		return nil, errors.New("not a node packet")
	}

	return &NodeTag{raw: b}, nil
}
