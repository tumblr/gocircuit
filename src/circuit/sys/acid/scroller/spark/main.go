package main

import (
	_ "circuit/load"
	"circuit/sys/acid/scroller/worker"
	"circuit/use/circuit"
)

func main() {
	_, addr, err := circuit.Spawn("localhost", []string{"/scroller"}, worker.App{})
	if err != nil {
		println("Oh oh", err.Error())
		return
	}
	println("Spawned", addr.String())
}
