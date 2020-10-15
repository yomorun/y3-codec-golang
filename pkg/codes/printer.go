package codes

import (
	"fmt"
	"io"
	"time"

	y3 "github.com/yomorun/yomo-codec-golang"
)

func PrintNodePacket(node *y3.NodePacket) {
	printNodeFormat(node, " %#X=%v ", false, true)
}

func printNodeFormat(node *y3.NodePacket, format string, isArray bool, isRoot bool) {
	if isRoot {
		fmt.Printf("%#x:{ ", node.SeqID())
	}

	if len(node.NodePackets) > 0 {
		for _, n := range node.NodePackets {
			if n.IsArray() {
				fmt.Printf(" %#x:[ ", n.SeqID())
				printNodeFormat(&n, format, true, false)
				fmt.Printf(" ]")
				continue
			}

			if n.SeqID() == 0x00 {
				fmt.Printf(" { ")
			} else {
				fmt.Printf(" %#x:{ ", n.SeqID())
			}
			printNodeFormat(&n, format, false, false)
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
		fmt.Printf(" }")
	}
}

type FmtOut struct{ io.Writer }

func (w FmtOut) Write(buf []byte) (int, error) {
	fmt.Printf("FmtOut: %s\n", FormatBytes(buf)) // debug:
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
