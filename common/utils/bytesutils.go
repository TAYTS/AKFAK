package utils

import (
	"encoding/binary"
	"io"
	"math"
)

// SizeOfVarint return the number of bytes needed to encode a varint (int)
func SizeOfVarint(value int) int {
	uValue := uint32(math.Abs(float64(value)))
	v := (uValue << 1) ^ (uValue >> 31)
	bytes := 1

	for (v & 0xffffff80) != 0 {
		bytes++
		v >>= 7
	}
	return bytes
}

// SizeOfVarintLong return the number of bytes needed to encode a varint (int64)
func SizeOfVarintLong(value int64) int {
	uValue := uint64(math.Abs(float64(value)))
	v := (uValue << 1) ^ (uValue >> 63)
	bytes := 1

	for (v & 0xffffffffffffff80) != 0 {
		bytes++
		v >>= 7
	}
	return bytes
}

// WriteVarint write a given int following the variable-length zig-zag encoding into the buffer
func WriteVarint(value int, out io.Writer) error {
	uValue := uint32(math.Abs(float64(value)))
	v := (uValue << 1) ^ (uValue >> 31)

	for (v & 0xffffff80) != 0 {
		buffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(buffer, (v&0x7f)|0x80)
		_, err := out.Write([]byte{buffer[0]})
		if err != nil {
			return err
		}
		v >>= 7
	}
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, (v&0x7f)|0x80)
	_, err := out.Write([]byte{buffer[0]})
	if err != nil {
		return err
	}

	return nil
}

// WriteVarlong write a given int64 following the variable-length zig-zag encoding into the buffer
func WriteVarlong(value int64, out io.Writer) error {
	uValue := uint64(math.Abs(float64(value)))
	v := (uValue << 1) ^ (uValue >> 63)

	for (v & 0xffffffffffffff80) != 0 {
		buffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(buffer, (v&0x7f)|0x80)
		_, err := out.Write([]byte{buffer[0]})
		if err != nil {
			return err
		}
		v >>= 7
	}
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, (v&0x7f)|0x80)
	_, err := out.Write([]byte{buffer[0]})
	if err != nil {
		return err
	}

	return nil
}
