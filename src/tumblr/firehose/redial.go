// Package firehose implements a connection to the Tumblr Firehose
package firehose

import (
	"sync"
	"time"
)

type RedialConn struct {
	sync.Mutex
	req       *Request
	conn      *Conn
	reLast    time.Time
	reSuccess int32
	reErr     int32
}

func Redial(req *Request) *RedialConn {
	rc := &RedialConn{req: req}
	//rc.conn, _ = Dial(rc.req)
	return rc
}

func (rc *RedialConn) Stat() (time.Time, int32, int32) {
	rc.Lock()
	defer rc.Unlock()
	return rc.reLast, rc.reSuccess, rc.reErr
}

func (rc *RedialConn) redial(onlyIfNil bool) {
	if !onlyIfNil {
		rc.conn = nil
	}
	var err error
	for rc.conn == nil {
		rc.reLast = time.Now()
		if rc.conn, err = Dial(rc.req); err == nil {
			rc.reSuccess++
			break
		}
		rc.reErr++
		println("Redial error:", err.Error())
		time.Sleep(time.Second)
	}
}

func (rc *RedialConn) Read() *Event {
	rc.Lock()
	defer rc.Unlock()

	var err error
	var ev *Event
	rc.redial(true)
	for {
		if ev, err = rc.conn.Read(); err != nil {
			rc.redial(false)
			continue
		}
		return ev
	}
	panic("u")
}

func (rc *RedialConn) ReadInterface(v interface{}) {
	rc.Lock()
	defer rc.Unlock()
	rc.redial(true)
	for {
		if err := rc.conn.ReadInterface(v); err != nil {
			rc.redial(false)
			continue
		}
		return
	}
	panic("u")
}

func (rc *RedialConn) ReadRaw() string {
	rc.Lock()
	defer rc.Unlock()
	var err error
	var raw string
	rc.redial(true)
	for {
		if raw, err = rc.conn.ReadRaw(); err != nil {
			rc.redial(false)
			continue
		}
		return raw
	}
	panic("u")
}

func (rc *RedialConn) Close() error {
	rc.Lock()
	defer rc.Unlock()
	var err error
	if rc.conn != nil {
		err = rc.conn.Close()
	}
	rc.conn = nil
	rc.req = nil
	return err
}
