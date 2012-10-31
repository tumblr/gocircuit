package transport

import (
	"testing"
)

func TestTransport(t *testing.T) {
	ch := make(chan int)
	t1, t2 := New(1, ":9000"), New(2, ":9001")

	go func() {
		c := t2.Accept()
		v, err := c.Read()
		if err != nil {
			t.Errorf("read (%s)", err)
		}
		if v.(int) != 3 {
			t.Errorf("value")
		}
		ch <- 1
	}()

	c12, err := t1.Dial(t2.Addr())
	if err != nil {
		t.Fatalf("dial, %s\n", err)
	}
	if err = c12.Write(int(3)); err != nil {
		t.Errorf("write (%s)", err)
	}
	<-ch
}
