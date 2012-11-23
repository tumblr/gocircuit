package firehose

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	freq := &Request{
		HostPort:      "",
		Username:      "",
		Password:      "",
		ApplicationID: "",
		ClientID:      "",
		Offset:        "",
	}

	conns := make([]*Conn, 2)
	for i, _ := range conns {
		fmt.Printf("dialing %d\n", i)
		var err error
		conns[i], err = Dial(freq)
		if err != nil {
			t.Fatalf("dial (%s)", err)
		}
	}
	fmt.Printf("reading\n")

	var lk sync.Mutex
	var nread int64
	var t0 time.Time = time.Now()

	for _, conn := range conns {
		go func(conn *Conn) {
			for {
				if _, err := conn.Read(); err != nil {
					t.Errorf("read (%s)", err)
					continue
				}
				lk.Lock()
				nread++
				k := nread
				lk.Unlock()
				if k%10 == 0 {
					fmt.Printf("%g read/sec\n", float64(k)/(float64(time.Now().Sub(t0))/1e9))
				}
			}
		}(conn)
	}
	<-(chan int)(nil)
}
