package utils

import "fmt"

// FormatBytes format bytes to string for print
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
