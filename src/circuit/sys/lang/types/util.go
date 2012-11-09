package types

import (
	"encoding/binary"
	"hash/fnv"
)

func sliceStringID32(sign []string) int32 {
	h := fnv.New32a()
	for _, s := range sign {
		h.Write([]byte(s))
	}
	return int32Bytes(h.Sum(nil))
}

func sliceStringID64(sign []string) int64 {
	h := fnv.New64a()
	for _, s := range sign {
		h.Write([]byte(s))
	}
	return int64Bytes(h.Sum(nil))
}

func int64Bytes(p []byte) int64 {
	return int64(binary.BigEndian.Uint64(p))
}

func int32Bytes(p []byte) int32 {
	return int32(binary.BigEndian.Uint32(p))
}
