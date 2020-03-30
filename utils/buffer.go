package utils

import "encoding/binary"

// BytesToInt convert byte array of size 4 to int32
func BytesToInt(data []byte) int {
	if len(data) != 4 {
		return 0
	}

	convertedData := int(binary.LittleEndian.Uint32(data))
	return convertedData
}

// IntToBytes convert int32 to byte array of size 4
func IntToBytes(data int) []byte {
	byteData := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteData, uint32(data))
	return byteData
}
