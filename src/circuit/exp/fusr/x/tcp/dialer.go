package tcp

import (
	"sync"
	"tumblr/balkan/x"
)

type Dialer struct {
	sync.Mutex
	open map[x.Addr]*link
}

func NewDialer() *Dialer {
	return &Dialer{
		open: make(map[x.Addr]*link),
	}
}

func (t *Dialer) Dial(addr x.Addr) x.Conn {
	t.Lock()
	l, ok := t.open[addr]
	if !ok {
		l = newDialLink(addr)
		t.open[addr] = l
	}
	t.Unlock()
	return l.Dial() // link.Dial may block
}
