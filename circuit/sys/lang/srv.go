package lang

import (
	"sync"
	"circuit/use/circuit"
)

type srvTabl struct {
	sync.Mutex
	name map[string]circuit.X
}

func (t *srvTabl) Init() *srvTabl {
	t.Lock()
	defer t.Unlock()
	t.name = make(map[string]circuit.X)
	return t
}

func (t *srvTabl) Add(name string, receiver interface{}) circuit.X {
	t.Lock()
	defer t.Unlock()
	if _, present := t.name[name]; present {
		panic("service already listening")
	}
	x := PermRef(receiver)
	t.name[name] = x
	return x
}

func (t *srvTabl) Get(name string) circuit.X {
	t.Lock()
	defer t.Unlock()
	return t.name[name]
}
