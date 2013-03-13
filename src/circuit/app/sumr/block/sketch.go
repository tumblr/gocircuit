package block

import (
	"time"
	"circuit/app/sumr"
)

// Sketch is the type of the value stored within the key-value store.
// This particular example implementation uses a sketch whose job is to
// be a simple counter. Users will typically make their own copy of the sumr
// service and substitute types like Sketch with ones suitable to their needs.
type Sketch struct {
	UpdateTime time.Time // Application-level timestamp of the key
	Key        sumr.Key
	Sum        float64
}
