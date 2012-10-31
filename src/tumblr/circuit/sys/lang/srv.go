package lang

import (
	"sync"
	"tumblr/circuit/use/lang"
)

type srvTabl struct {
	sync.Mutex
	name map[string]lang.X
}

func (t *srvTabl) Init() *srvTabl {
	t.Lock()
	defer t.Unlock()
	t.name = make(map[string]lang.X)
	return t
}

func (t *srvTabl) Add(name string, receiver interface{}) lang.X {
	t.Lock()
	defer t.Unlock()
	if _, present := t.name[name]; present {
		panic("service already listening")
	}
	x := PermRef(receiver)
	t.name[name] = x
	return x
}

func (t *srvTabl) Get(name string) lang.X {
	t.Lock()
	defer t.Unlock()
	return t.name[name]
}
