package lang

import (
	"bytes"
	"circuit/use/circuit"
	"encoding/gob"
	"io"
	"sync"
)

func NewBytesConn(addr string) circuit.Conn {
	var b bytes.Buffer
	return ReadWriterConn(stringAddr(addr), nopCloser{&b})
}

type nopCloser struct {
	io.ReadWriter
}

func (nc nopCloser) Close() error {
	return nil
}

type stringAddr string

func (a stringAddr) Host() string {
	panic("not network host underlying")
}

func (a stringAddr) WorkerID() circuit.WorkerID {
	return 0
}

func (a stringAddr) String() string {
	return string(a)
}

// ReadWriterConn converts an io.ReadWriteClosert into a Conn
func ReadWriterConn(addr circuit.Addr, rwc io.ReadWriteCloser) circuit.Conn {
	return &readWriterConn{
		addr: addr,
		rwc:  rwc,
		enc:  gob.NewEncoder(rwc),
		dec:  gob.NewDecoder(rwc),
	}
}

type readWriterConn struct {
	addr circuit.Addr
	sync.Mutex
	rwc io.ReadWriteCloser
	enc *gob.Encoder
	dec *gob.Decoder
}

type blob struct {
	Cargo interface{}
}

func (conn *readWriterConn) Read() (interface{}, error) {
	conn.Lock()
	defer conn.Unlock()
	var b blob
	err := conn.dec.Decode(&b)
	if err != nil {
		return nil, err
	}
	return b.Cargo, nil
}

func (conn *readWriterConn) Write(cargo interface{}) error {
	conn.Lock()
	defer conn.Unlock()
	return conn.enc.Encode(&blob{cargo})
}

func (conn *readWriterConn) Close() error {
	conn.Lock()
	defer conn.Unlock()
	return conn.rwc.Close()
}

func (conn *readWriterConn) Addr() circuit.Addr {
	return conn.addr
}
