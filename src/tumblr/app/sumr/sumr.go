package sumr

import (
	"fmt"
)

// Key is used for hashing features
type Key int64

func (k Key) String() string {
	return fmt.Sprintf("K%016x", int64(k))
}
