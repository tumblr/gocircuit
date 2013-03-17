// Package iomisc implements miscellaneous I/O facilities
package iomisc

import (
	"io"
	"sync"
)

type combinedReader struct {
	pipe   *io.PipeReader
	wlk    sync.Mutex
	closed int
}

// Combine returns an io.Reader that greedily reads from r1 and r2 in parallel
func Combine(r1, r2 io.Reader) io.Reader {
	pr, pw := io.Pipe()
	c := &combinedReader{pipe: pr}
	go c.readTo(r1, pw)
	go c.readTo(r2, pw)
	return c
}

func (c *combinedReader) readTo(r io.Reader, w *io.PipeWriter) {
	p := make([]byte, 1e5)
	for {
		n, err := r.Read(p)
		if n > 0 {
			c.wlk.Lock()
			w.Write(p[:n])
			c.wlk.Unlock()
		}
		if err != nil {
			c.wlk.Lock()
			defer c.wlk.Unlock()
			c.closed++
			if c.closed == 2 {
				w.Close()
			}
			return
		}
	}
}

func (c *combinedReader) Read(p []byte) (int, error) {
	return c.pipe.Read(p)
}
