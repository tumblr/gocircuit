package main

import (
	"tumblr/circuit/use/lang"
	"tumblr/circuit/sys/n/trojan"
	"tumblr/circuit/sys/transport"
)

func NewTransport(id lang.RuntimeID, addr string) lang.Transport {
	return transport.New(id, addr)
}

func main() {
	trojan.Main(NewTransport)
}
