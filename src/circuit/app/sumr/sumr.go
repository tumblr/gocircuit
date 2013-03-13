// Package sumr implements a circuit application for a distributed, persistent key-value store
package sumr

import "fmt"

// Key is the type of keys used in the sumr database
type Key int64

// String returns the textual representation of this key
func (k Key) String() string {
	return fmt.Sprintf("K%016x", int64(k))
}
