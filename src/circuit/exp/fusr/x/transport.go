package x

import (
	"strconv"
	"strings"
)

type Addr string

func (addr Addr) Port() int {
	i := strings.Index(string(addr), ":")
	if i < 0 {
		panic("endpoint address has no port")
	}
	port, err := strconv.Atoi(string(addr)[i+1:])
	if err != nil {
		panic("endpoint address with invalid port")
	}
	return port
}

type Conn interface {
	Read() (interface{}, error)
	Write(interface{}) error
	Close() error
	RemoteAddr() Addr
}

type Listener interface {
	Accept() Conn
}

type Dialer interface {
	Dial(addr Addr) Conn
}

type Transport interface {
	Dialer
	Listener
}
