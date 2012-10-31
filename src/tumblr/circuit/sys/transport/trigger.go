package transport

import (
	"sync"
)

type Trigger struct {
	lk       sync.Mutex
	engaged  bool
	nwaiters int
	ch       chan struct{}
}

func (t *Trigger) Lock() bool {
	t.lk.Lock()
	if t.ch == nil {
		t.ch = make(chan struct{})
	}
	if t.engaged {
		t.nwaiters++
		t.lk.Unlock()
		<-t.ch
		return false
	}
	t.engaged = true
	t.lk.Unlock()
	return true
}

func (t *Trigger) Unlock() {
	t.lk.Lock()
	defer t.lk.Unlock()
	if !t.engaged {
		panic("unlocking a non-engaged trigger")
	}
	for t.nwaiters > 0 {
		t.ch <- struct{}{}
		t.nwaiters--
	}
	t.engaged = false
}
