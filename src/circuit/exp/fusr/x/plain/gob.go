package plain

import (
	"encoding/gob"
	"io"
	"tumblr/balkan/x"
)

type ReadWriter interface {
	Read() (interface{}, error)
	Write(interface{}) error
}

type Closer interface {
	Close() error
}

type ReadWriteCloser interface {
	ReadWriter
	Closer
}

// gobConn implements ReadWriteCloser on top of a io.ReadWriteCloser
type gobConn struct {
	*gob.Encoder
	*gob.Decoder
	io.ReadWriteCloser
}

type gobMsg struct {
	Payload interface{}
}

func newGobConn(c io.ReadWriteCloser) *gobConn {
	return &gobConn{
		Encoder:         gob.NewEncoder(c),
		Decoder:         gob.NewDecoder(c),
		ReadWriteCloser: c,
	}
}

func (g *gobConn) Read() (interface{}, error) {
	var msg gobMsg
	if err := g.Decode(&msg); err != nil {
		return nil, err
	}
	return msg.Payload, nil
}

func (g *gobConn) Write(v interface{}) error {
	var msg gobMsg = gobMsg{v}
	return g.Encode(&msg)
}

func (g *gobConn) RemoteAddr() x.Addr {
	return x.Addr("hello")
}
