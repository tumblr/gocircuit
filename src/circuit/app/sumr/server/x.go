package server

import (
	"circuit/use/circuit"
	"time"
)

// Main wraps the worker function that starts a sumr shard server
type main struct{}

func init() {
	circuit.RegisterFunc(main{})
}

// Main starts a sumr shard server
// diskpath is a directory path on the local file system, where the function is executed,
// where the shard will persist its data.
func (main) Main(diskpath string, forgetafter time.Duration) (circuit.XPerm, error) {
	srv, err := New(diskpath, forgetafter)
	if err != nil {
		return nil, err
	}
	circuit.Daemonize(func() { <-(chan int)(nil) })
	return circuit.PermRef(srv), nil
}
