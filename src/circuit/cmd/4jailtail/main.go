package main

import (
	"fmt"
	"os"
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage:", os.Args[0], "AnchorPath PathWithinJail")
		os.Exit(1)
	}
	file, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)", err)
		os.Exit(1)
	}
	x, err := circuit.TryDial(file.Owner(), "acid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem dialing 'acid' service (%s)", err)
		os.Exit(1)
	}
	println("Connected")

	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "Worker disappeared during call (%#v)", p)
			os.Exit(1)
		}
	}()

	retrn := x.Call("JailOpen", os.Args[2])
	println("Got it")
	_/*fclix*/, err = retrn[0].(circuit.X), retrn[1].(error)
	if err != nil {
		println("jail open error", err.Error())
		os.Exit(1)
	}
}
