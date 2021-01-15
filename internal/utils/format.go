package utils

import "fmt"

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
