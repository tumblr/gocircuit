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

// String returns a cononical string representation of this worker ID.
func (r WorkerID) String() string {
	return fmt.Sprintf("R%016x", int64(r))
}

// ChooseWorkerID returns a random worker ID.
func ChooseWorkerID() WorkerID {
	return WorkerID(rand.Int63())
}

// ParseOrHashWorkerID tries to parse the string s as a canonical worker ID representation.
// If it fails, it treats s as an unconstrained string and hashes it to a worker ID value.
// In either case, it returns a WorkerID value.
func ParseOrHashWorkerID(s string) WorkerID {
	id, err := ParseWorkerID(s)
	if err != nil {
		return HashWorkerID(s)
	}
	return id
}

// ParseWorkerID parses the string s for a canonical representation of a worker
// ID and returns a corresponding WorkerID value.
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

// HashWorkerID hashes the unconstrained string s into a worker ID value.
func HashWorkerID(s string) WorkerID {
	h := fnv.New64a()
	h.Write([]byte(s))
	return WorkerID(h.Sum64())
}
