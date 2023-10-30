package helper

import (
	"fmt"
)

/*
	Formatting []byte with %x gives a string of hex
	Want it still in structure of a []byte
	Returns a string representative of a []byre encoded in hex
*/
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