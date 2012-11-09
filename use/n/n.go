package n

import (
	"encoding/gob"
	"io"
	"circuit/kit/join"
	"circuit/use/circuit"
)

// Host
type Host struct {
	Host string
}
func init() {
	gob.Register(&Host{})
}

func ParseHost(host string) circuit.Host {
	return &Host{host}
}

func (h Host) String() string {
	return h.Host
}

// Process
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
