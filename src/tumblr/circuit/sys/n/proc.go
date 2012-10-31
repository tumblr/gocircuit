package n

import (
	"io"
	"tumblr/circuit/use/lang"
	"tumblr/circuit/sys/transport"
)

type Console struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

type Process struct {
	console Console
	addr    *transport.Addr
}

func (p *Process) Addr() lang.Addr {
	return p.addr
}

func (p *Process) Kill() error {
	return kill(p.addr)
}

func (p *Process) Stdin()  io.WriteCloser {
	panic("ni")
	return p.console.stdin
}

func (p *Process) Stdout() io.ReadCloser {
	panic("ni")
	return p.console.stdout
}

func (p *Process) Stderr() io.ReadCloser {
	panic("ni")
	return p.console.stderr
}
