package server

import (
	"time"
	"circuit/use/circuit"
)

// Main wraps the worker function that starts a sumr shard server
type Main struct{}

func init() {
	circuit.RegisterFunc(Main{})
}

// Main starts a sumr shard server
// diskpath is a directory path on the local file system, where the function is executed,
// where the shard will persist its data.
func (Main) Main(diskpath string, forgetafter time.Duration) (circuit.XPerm, error) {
	srv, err := New(diskpath, forgetafter)
	if err != nil {
		return nil, err
	}
	circuit.Daemonize(func() { <-(chan int)(nil) })
	return circuit.PermRef(srv), nil
}
