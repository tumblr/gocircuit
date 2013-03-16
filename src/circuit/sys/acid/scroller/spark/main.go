package main

import (
	_ "circuit/load"
	"circuit/use/circuit"
	"circuit/sys/acid/scroller/worker"
)

func main() {
	_, addr, err := circuit.Spawn("localhost", []string{"/scroller"}, worker.App{})
	if err != nil {
		println("Oh oh", err.Error())
		return
	}
	println("Spawned", addr.String())
}
