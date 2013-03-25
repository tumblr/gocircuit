// Package join provides a mechanism for linking an implementation package to a declaration package
package join

import (
	"sync"
)

// SetThenGet is a synchronized interface value, which can be set once and read many times
type SetThenGet struct {
	Name string
	lk   sync.Mutex
	v    interface{}
}

// Set sets the value to v
func (j *SetThenGet) Set(v interface{}) {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v != nil {
		panic(j.Name + " already set")
	}
	j.v = v
}

// Get returns this value
func (j *SetThenGet) Get() interface{} {
	j.lk.Lock()
	defer j.lk.Unlock()
	if j.v == nil {
		panic(j.Name + " not set")
	}
	return j.v
}
