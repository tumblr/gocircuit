package transport

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"tumblr/circuit/use/lang"
)

// gobConn keeps a Conn instance together with its gob codecs
type gobConn struct {
	*gob.Encoder
	*gob.Decoder
	net.Conn
}

func newGobConn(c net.Conn) *gobConn {
	return &gobConn{
		Encoder:  gob.NewEncoder(c),
		Decoder:  gob.NewDecoder(c),
		Conn:     c,
	}
}

// Transport ..
// Transport implements lang.Transport, lang.Dialer and lang.Listener
type Transport struct {
	self       lang.Addr
	bind       *Addr
	listener   *net.TCPListener
	addrtabl   *addrTabl
	// How many unacknowledged messages we are willing to keep per link, before
	// we start blocking on writes
	pipelining int  

	lk         sync.Mutex
	remote     map[lang.RuntimeID]*link

	ach        chan *conn
}

func NewClient(id lang.RuntimeID) *Transport {
	return New(id, "", "localhost")
}

const DefaultPipelining = 333

func New(id lang.RuntimeID, bindAddr string, host string) *Transport {

	// Bind
	var l *net.TCPListener
	if strings.Index(bindAddr, ":") < 0 {
		bindAddr = bindAddr + ":0"
	}
	l_, err := net.Listen("tcp", bindAddr)
	if err != nil {
		panic(err)
	}
	
	// Build transport structure
	l = l_.(*net.TCPListener)
	t := &Transport{
		listener:   l,
		addrtabl:   makeAddrTabl(),
		pipelining: DefaultPipelining,
		remote:     make(map[lang.RuntimeID]*link),
		ach:        make(chan *conn),
	}

	// Resolve self address
	laddr := l.Addr().(*net.TCPAddr)
	t.self, err = NewAddr(id, os.Getpid(), fmt.Sprintf("%s:%d", host, laddr.Port))
	if err != nil {
		panic(err)
	}

	// This LocalAddr might be useless for connect purposes (e.g. 0.0.0.0). Consider self instead.
	t.bind = t.addrtabl.Normalize(&Addr{ID: id, PID: os.Getpid(), Addr: laddr})

	go t.loop()
	return t
}

func (t *Transport) Port() int {
	return t.bind.Addr.Port
}

func (t *Transport) Addr() lang.Addr {
	return t.self
}

func (t *Transport) Accept() lang.Conn {
	return <-t.ach
}

func (t *Transport) loop() {
	for {
		c, err := t.listener.AcceptTCP()
		if err != nil {
			panic(err)  // Best not to be quiet about it
		}
		t.link(c)
	}
}

func (t *Transport) Dial(a lang.Addr) (lang.Conn, error) {
	a_ := a.(*Addr)
	t.lk.Lock()
	l, ok := t.remote[a_.ID]
	t.lk.Unlock()
	if ok {
		return l.Open()
	}
	l, err := t.dialLink(a_)
	if err != nil {
		return nil, err
	}
	return l.Open()
}

func (t *Transport) dialLink(a *Addr) (*link, error) {
	a = t.addrtabl.Normalize(a)
	c, err := net.DialTCP("tcp", nil, a.Addr)
	if err != nil {
		return nil, err
	}
	l, err := t.link(c)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (t *Transport) drop(id lang.RuntimeID) {
	t.lk.Lock()
	delete(t.remote, id)
	t.lk.Unlock()
}

func (t *Transport) link(c *net.TCPConn) (*link, error) {
	g := newGobConn(c)

	// Send-receive welcome, ala mutual authentication
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		g.Encode(&welcomeMsg{ID: t.bind.ID, PID: os.Getpid()})
		wg.Done()
	}()
	var welcome welcomeMsg
	if err := g.Decode(&welcome); err != nil {
		wg.Wait() // Wait to finish sending, so no compete
		g.Close()
		return nil, err
	}
	wg.Wait() // Wait to finish sending welcome msg

	addr := t.addrtabl.Normalize(&Addr{
		ID:   welcome.ID, 
		PID:  welcome.PID,
		Addr: c.RemoteAddr().(*net.TCPAddr),
	})

	t.lk.Lock()
	l, ok := t.remote[addr.ID]
	if !ok {
		l = makeLink(addr, g, t.ach, func() { t.drop(addr.ID) }, t.pipelining)
		t.remote[addr.ID] = l
		t.lk.Unlock()
	} else {
		t.lk.Unlock()
		if err := l.acceptReconnect(g); err != nil {
			g.Close()
			return nil, err
		}
	}
	return l, nil
}

func (t *Transport) Dialer() lang.Dialer { 
	return t 
}

func (t *Transport) Listener() lang.Listener { 
	return t 
}

func (t *Transport) Close() {
	panic("not supported")
}
