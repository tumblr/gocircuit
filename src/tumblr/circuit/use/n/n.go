package n

import (
	"io"
	"tumblr/circuit/kit/join"
	"tumblr/circuit/use/lang"
)

type Process interface {
	Console
	Addr() lang.Addr
	Kill() error
}

type Console interface {
	Stdin()  io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}

// Process is a usually remote OS process of a circuit runtime.

func Spawn(host lang.Host, anchors ...string) (Process, error) {
	return get().Spawn(host, anchors...)
}

func Kill(addr lang.Addr) error {
	return get().Kill(addr)
}

type Commander interface {
	Spawn(lang.Host, ...string) (Process, error)
	Kill(lang.Addr) error
}

// Binding mechanism
var link = join.SetThenGet{Name: "commander system"}

func Bind(v Commander) {
	link.Set(v)
}

func get() Commander {
	return link.Get().(Commander)
}
