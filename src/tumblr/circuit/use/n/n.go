package n

import (
	"io"
	"tumblr/circuit/kit/join"
	"tumblr/circuit/use/circuit"
)

type Process interface {
	Console
	Addr() circuit.Addr
	Kill() error
}

type Console interface {
	Stdin()  io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}

// Process is a usually remote OS process of a circuit runtime.

func Spawn(host circuit.Host, anchors ...string) (Process, error) {
	return get().Spawn(host, anchors...)
}

func Kill(addr circuit.Addr) error {
	return get().Kill(addr)
}

type Commander interface {
	Spawn(circuit.Host, ...string) (Process, error)
	Kill(circuit.Addr) error
}

// Binding mechanism
var link = join.SetThenGet{Name: "commander system"}

func Bind(v Commander) {
	link.Set(v)
}

func get() Commander {
	return link.Get().(Commander)
}
