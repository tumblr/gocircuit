package lang

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"
)

var ErrParse = NewError("parse")

// RuntimeID ...
type RuntimeID uint64

func (r RuntimeID) String() string {
	return fmt.Sprintf("R%016x", int64(r))
}

// ChooseRuntimeID returns a random runtime ID
func ChooseRuntimeID() RuntimeID {
	return RuntimeID(rand.Int63())
}

func ParseOrHashRuntimeID(s string) RuntimeID {
	id, err := ParseRuntimeID(s)
	if err != nil {
		return HashRuntimeID(s)
	}
	return id
}

func ParseRuntimeID(s string) (RuntimeID, error) {
	if len(s) != 17 || s[0] != 'R' {
		return 0, ErrParse
	}
	ui64, err := strconv.ParseInt(s[1:], 16, 64)
	if err != nil {
		return 0, ErrParse
	}
	return RuntimeID(ui64), nil
}

func HashRuntimeID(s string) RuntimeID {
	h := fnv.New64a()
	h.Write([]byte(s))
	return RuntimeID(h.Sum64())
}
