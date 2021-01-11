package packetutils

import (
	"fmt"
	"io"
	"time"

	y3 "github.com/yomorun/y3-codec-golang"
)

func PrintNodePacket(node *y3.NodePacket) {
	PrintNodeFormat(node, " %#X=%v ", false, true)
}

func PrintArrayPacket(node *y3.NodePacket) {
	//Parsing [0xc1, 0x06, 0x00, 0x01, 0x61, 0x00, 0x01, 0x62] EQUALS 0xc1:[0x02,0x04]")
	PrintNodeFormat(node, " %#X=%v ", true, true)
}

func PrintNodeFormat(node *y3.NodePacket, format string, isArray bool, isRoot bool) {
	if isRoot {
		if isArray {
			fmt.Printf("%#x:[ ", node.SeqID())
		} else {
			fmt.Printf("%#x:{ ", node.SeqID())
		}
	}

	if len(node.NodePackets) > 0 {
		for _, n := range node.NodePackets {
			if n.IsArray() {
				fmt.Printf(" %#x:[ ", n.SeqID())
				PrintNodeFormat(&n, format, true, false)
				fmt.Printf(" ]")
				continue
			}

			if n.SeqID() == 0x00 {
				fmt.Printf(" { ")
			} else {
				fmt.Printf(" %#x:{ ", n.SeqID())
			}
			PrintNodeFormat(&n, format, false, false)
			fmt.Printf(" }")

		}
	}
	if len(node.PrimitivePackets) > 0 {
		for _, p := range node.PrimitivePackets {
			if isArray {
				fmt.Printf(" %#x ", p.ToBytes())
			} else {
				fmt.Printf(format, p.SeqID(), fmt.Sprintf("%#x", p.ToBytes()))
			}
		}
	}

	if isRoot {
		if isArray {
			fmt.Printf(" ]")
		} else {
			fmt.Printf(" }")
		}
	}
}

type FmtOut struct{ io.Writer }

func (w FmtOut) Write(buf []byte) (int, error) {
	res, _, _ := y3.DecodeNodePacket(buf)
	fmt.Printf("%v:\t", time.Now().Format("2006-01-02 15:04:05")) // debug:
	PrintNodePacket(res)
	fmt.Println()
	return 0, nil
}

func FormatBytes(buf []byte) string {
	var str = ""
	for i, c := range buf {
		if i == 0 {
			str = str + fmt.Sprintf("%#x", c)
			continue
		}
		str = str + fmt.Sprintf(" %#x", c)
	}
	return str
}
