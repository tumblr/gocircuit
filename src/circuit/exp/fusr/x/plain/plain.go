package plain

import (
	"net"
	"time"
	"tumblr/balkan/x"
)

// Listener
type Listener struct {
	l net.Listener
}

func NewListener(addr x.Addr) *Listener {
	l, err := net.Listen("tcp", string(addr))
	if err != nil {
		panic(err)
	}
	return &Listener{l}
}

func (l *Listener) Accept() x.Conn {
	c, err := l.l.Accept()
	if err != nil {
		panic(err)
	}
	return newGobConn(c)
}

// Dialer
type Dialer struct{}

func NewDialer() Dialer {
	return Dialer{}
}

func (Dialer) Dial(addr x.Addr) x.Conn {
	for {
		tcpaddr, err := net.ResolveTCPAddr("tcp", string(addr))
		if err != nil {
			println("tcp resolve:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		c, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			println("tcp dial:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		return newGobConn(c)
	}
	panic("u")
}
