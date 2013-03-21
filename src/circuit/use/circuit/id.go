package circuit

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strconv"
)

var ErrParse = NewError("parse")

// WorkerID represents the ID of a circuit worker process.
type WorkerID uint64

func (r WorkerID) String() string {
	return fmt.Sprintf("R%016x", int64(r))
}

// ChooseWorkerID returns a random runtime ID
func ChooseWorkerID() WorkerID {
	return WorkerID(rand.Int63())
}

func ParseOrHashWorkerID(s string) WorkerID {
	id, err := ParseWorkerID(s)
	if err != nil {
		return HashWorkerID(s)
	}
	return id
}

func ParseWorkerID(s string) (WorkerID, error) {
	if len(s) != 17 || s[0] != 'R' {
		return 0, ErrParse
	}
	ui64, err := strconv.ParseInt(s[1:], 16, 64)
	if err != nil {
		return 0, ErrParse
	}
	return WorkerID(ui64), nil
}

func HashWorkerID(s string) WorkerID {
	h := fnv.New64a()
	h.Write([]byte(s))
	return WorkerID(h.Sum64())
}
