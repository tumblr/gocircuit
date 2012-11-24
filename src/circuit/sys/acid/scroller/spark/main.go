package main

import (
	_ "circuit/load"
	"circuit/use/circuit"
	"circuit/use/n"
	"circuit/sys/acid/scroller/worker"
)

func main() {
	_, addr, err := circuit.Spawn(n.ParseHost("localhost"), []string{"/scroller"}, worker.App{})
	if err != nil {
		println("Oh oh", err.Error())
		return
	}
	println("Spawned", addr.String())
}
