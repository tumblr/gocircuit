package main

import (
	"circuit/exp/telefile/srv"
	"circuit/kit/tele/file"

	_ "circuit/load"
	"circuit/use/circuit"
	"io"
	"os"
)

func main() {
	println("Starting")
	r, _, err := circuit.Spawn("localhost", []string{"/telefile"}, srv.App{}, "/tmp/telehelo")
	if err != nil {
		println("Oh oh", err.Error())
		return
	}
	fcli := file.NewFileClient(r[0].(circuit.X))
	defer func() {
		recover()
	}()
	io.Copy(os.Stdout, fcli)
}
