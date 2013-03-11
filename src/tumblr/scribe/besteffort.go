package scribe

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

func BestEffortDial(hostport string) (*BestEffortConn, error) {
	be := &BestEffortConn{
		conn:     nil,
		hostport: hostport,
	}
	go be.redial()
	return be, nil
}

var ErrRedialing = errors.New("redialing")

func (bec *BestEffortConn) E(category, payload string) error {
	bec.Lock()
	defer bec.Unlock()
	if bec.conn == nil {
		return ErrRedialing
	}

	err := bec.conn.E(category, payload)
	if err != nil {
		// If we are dealing with a network error, than spawn a redial
		bec.conn = nil
		go bec.redial()
	}
	return err
}

func (bec *BestEffortConn) Emit(msgs ...Message) error {
	bec.Lock()
	defer bec.Unlock()

	if bec.conn == nil {
		return ErrRedialing
	}

	err := bec.conn.Emit(msgs...)
	if err != nil {
		// If we are dealing with a network error, than spawn a redial
		bec.conn = nil
		go bec.redial()
	}
	return err
}

func (bec *BestEffortConn) redial() {
	var err error
	var reconn *Conn
	for reconn == nil {
		time.Sleep(2*time.Second)  // Sleep a bit so things don't spin out of control
		bec.Lock()
		hostport := bec.hostport
		bec.Unlock()
		if hostport == "" {
			// The BestEffortConn has been closed
			return
		}
		if reconn, err = Dial(hostport); err != nil {
			reconn = nil
		}
	}

	bec.Lock()
	defer bec.Unlock()
	if bec.hostport == "" {
		reconn.Close()
	} else {
		bec.conn = reconn
	}
}

func (bec *BestEffortConn) Close() error {
	bec.Lock()
	defer bec.Unlock()
	bec.hostport = ""
	if bec.conn != nil {
		bec.conn.Close()
	}
	return nil
}
