package helper

import (
	"fmt"
)

func FormatByteSliceAsHexSliceString(bytes []byte) string {
	out := "["
	for i, b := range bytes {
		out += fmt.Sprintf("%x", b)
		if i != len(bytes) - 1 {
			out += " "
		}
	}
	out += "]"
	return out
}