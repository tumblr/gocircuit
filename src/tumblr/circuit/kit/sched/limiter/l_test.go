package limiter

import (
	"time"
	"testing"
)

func TestLimiter(t *testing.T) {
	l := New(2)
	for i := 0; i < 9; i++ {
		i_ := i
		l.Go(func() { 
			println("{", i_) 
			time.Sleep(time.Second); 
			println("}", i_) 
		})
	}
	l.Wait()
	// TODO: Test that all routines open and close
	println("DONE")
}
