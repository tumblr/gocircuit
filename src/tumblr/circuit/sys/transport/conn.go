package transport

import (
	"math/rand"
	"sync"
	"tumblr/circuit/use/lang"
)

// Within a TCP connection, the connID distinguishes a unique logical session
type connID int32

func chooseConnID() connID {
	return connID(rand.Int31())
}

// conn implements lang.Conn
type conn struct {
	id   connID
	addr *Addr
	ann  bool // If ann(ounced) is true, we need not set the First flag on outgoing messages

	lk   sync.Mutex // conn.Close and link.readLoop are competing for send/close to ch
	ch   chan interface{} // link.readLoop send msgs for this conn to conn.Read
	l    *link
}

func makeConn(id connID, l *link) *conn {
	return &conn{id: id, addr: l.addr, l: l, ch: make(chan interface{})}
}

func (c *conn) Read() (interface{}, error) {
	v, ok := <-c.ch
	if !ok {
		return nil, ErrEnd
	}
	return v, nil
}

func (c *conn) sendRead(v interface{}) {
	c.lk.Lock() // Lock ch to send payload to it
	if c.l != nil { // Implies c.ch not closed
		c.ch <- v
	}
	c.lk.Unlock()
}

func (c *conn) Write(v interface{}) error {
	msg := &connMsg{ID: c.id, Payload: v}
	c.lk.Lock()
	l := c.l
	c.lk.Unlock()
	if l == nil {
		return ErrEnd
	}
	return l.Write(msg)
}

// Close instructs link to remove it from the list of open connections.
// For efficiency Close does not send any network messages.
// Users must ensure they close explicitly.
func (c *conn) Close() error {
	c.lk.Lock()
	defer c.lk.Unlock()

	if c.l == nil {
		return nil
	}
	close(c.ch)
	c.l.drop(c.id)
	c.l = nil
	return nil
}

func (c *conn) Addr() lang.Addr {
	return c.addr
}
