package tcp

import (
	"net/http"
	"testing"
	x "circuit/exp/teleport"
	_ "circuit/kit/debug/http/trace"
	"sync"
)

func init() {
	go http.ListenAndServe(":1505", nil)
}

func TestTransport(t *testing.T) {
	const N = 100
	ch := make(chan int)
	d := NewDialer()
	laddr := x.Addr("localhost:9001")
	l := NewListener(":9001")

	go func() {
		for i := 0; i < N; i++ {
			c := l.Accept()
			v, err := c.Read()
			if err != nil {
				t.Errorf("read (%s)", err)
			}
			if v.(int) != 3 {
				t.Errorf("value")
			}
			c.Close()
		}
		ch <- 1
	}()

	var slk sync.Mutex
	sent := 0
	for i := 0; i < N; i++ {
		go func() {
			c := d.Dial(laddr)
			if err := c.Write(int(3)); err != nil {
				t.Errorf("write (%s)", err)
			}
			c.Close()
			slk.Lock()
			defer slk.Unlock()
			sent++
		}()
	}
	<-ch
}
