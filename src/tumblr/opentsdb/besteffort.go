package opentsdb

import (
	"errors"
	"sync"
	"time"
)

type BestEffortConn struct {
	sync.Mutex
	conn *Conn
	hostport string
}

var ErrRedialing = errors.New("redialing")

func BestEffortDial(hostport string) (*BestEffortConn, error) {
	be := &BestEffortConn{
		hostport: hostport,
		conn:     nil,
	}
	go be.redial()
	return be, nil
}

func (c *BestEffortConn) Put(metric string, value interface{}, tags ...Tag) error {
	c.Lock()
	defer c.Unlock()
	if c.conn == nil {
		return ErrRedialing
	}

	err := c.conn.Put(metric, value, tags...)
	if err != nil && err != ErrArg {
		// If we are dealing with a network error, than spawn a redial
		c.conn = nil
		go c.redial()
	}
	return err
}

func (c *BestEffortConn) redial() {
	var err error
	var conn *Conn
	for conn == nil {
		time.Sleep(2*time.Second)  // Sleep a bit so things don't spin out of control
		c.Lock()
		hostport := c.hostport
		c.Unlock()
		if hostport == "" {
			return
		}
		if conn, err = Dial(hostport); err != nil {
			conn = nil
		}
	}

	c.Lock()
	defer c.Unlock()
	if c.hostport == "" {
		conn.Close()
	} else {
		c.conn = conn
	}
}

func (c *BestEffortConn) Close() error {
	c.Lock()
	defer c.Unlock()
	c.hostport = ""
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
