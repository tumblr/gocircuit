package kafka

import (
	"encoding/binary"
)

// 64-bit

func Int64Bytes(value int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(value))
	return result
}

func Uint64Bytes(value uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, value)
	return result
}

// 32-bit

func Int32Bytes(value int32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, uint32(value))
	return result
}

func Uint32Bytes(value uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, value)
	return result
}

// 16-bit

func Int16Bytes(value int16) []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(value))
	return result
}

// 64-bit 

func BytesInt64(p []byte) int64 {
	return int64(binary.BigEndian.Uint64(p))
}

func BytesUint64(p []byte) uint64 {
	return binary.BigEndian.Uint64(p)
}

// 32-bit 

func BytesInt32(p []byte) int32 {
	return int32(binary.BigEndian.Uint32(p))
}

func BytesUint32(p []byte) uint32 {
	return binary.BigEndian.Uint32(p)
}

// 16-bit

func BytesInt16(p []byte) int16 {
	return int16(binary.BigEndian.Uint16(p))
}
