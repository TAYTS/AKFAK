package utils

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// StringToBytes convert string to bytearray
func StringToBytes(str string) []byte {
	return []byte(str)
}

// UTF8Length return the byte length of utf-8 string
func UTF8Length(str string) int {
	count := 0

	for _, char := range str {
		count += utf8.RuneLen(char)
	}

	return count
}

// WriteTo write the contents of buffer to an output stream. The bytes are copied from the current position in the buffer.
func WriteTo(out io.Writer, buffer *bytes.Reader) error {
	_, err := buffer.WriteTo(out)
	return err
}
