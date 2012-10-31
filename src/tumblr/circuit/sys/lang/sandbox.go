package lang

import (
	"encoding/gob"
	"io"
	"sync"
	"tumblr/circuit/use/lang"
)

type sandbox struct {
	lk  sync.Mutex
	l   map[lang.RuntimeID]*listener
}

var s = &sandbox{l: make(map[lang.RuntimeID]*listener)}

// NewSandbox creates a new transport instance, part of a sandbox network in memory
func NewSandbox() lang.Transport {
	s.lk.Lock()
	defer s.lk.Unlock()

	l := &listener{
		id: lang.ChooseRuntimeID(),
		ch: make(chan *halfconn),
	}
	l.a = &addr{ID: l.id, l: l}
	s.l[l.id] = l
	return l
}

func dial(remote lang.Addr) (lang.Conn, error) {
	pr, pw := io.Pipe()
	qr, qw := io.Pipe()
	srvhalf := &halfconn{PipeWriter: qw, PipeReader: pr}
	clihalf := &halfconn{PipeWriter: pw, PipeReader: qr}
	s.lk.Lock()
	l := s.l[remote.(*addr).RuntimeID()]
	s.lk.Unlock()
	if l == nil {
		panic("unknown listener id")
	}
	go func() { 
		l.ch <- srvhalf
	}()
	return ReadWriterConn(l.Addr(), clihalf), nil
}

// addr implements Addr
type addr struct {
	ID  lang.RuntimeID
	l   *listener
}

func (a *addr) RuntimeID() lang.RuntimeID {
	return a.ID
}

func (a *addr) String() string {
	return a.ID.String()
}

func init() {
	gob.Register(&addr{})
}

// listener implements Listener
type listener struct {
	id  lang.RuntimeID
	a   *addr
	ch  chan *halfconn
}

func (l *listener) Addr() lang.Addr { 
	return l.a
}

func (l *listener) Accept() lang.Conn { 
	return ReadWriterConn(l.Addr(), <-l.ch)
}

func (l *listener) Close() {
	s.lk.Lock()
	defer s.lk.Unlock()
	delete(s.l, l.id)
}

func (l *listener) Dial(remote lang.Addr) (lang.Conn, error) {
	return dial(remote)
}

// halfconn is one end of a byte-level connection
type halfconn struct {
	*io.PipeReader
	*io.PipeWriter
}

func (h *halfconn) Close() error {
	return h.PipeWriter.Close()
}
