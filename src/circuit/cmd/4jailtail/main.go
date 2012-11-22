package main

import (
	"fmt"
	"os"
	_ "circuit/load"
	"circuit/kit/tele/file"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"io"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage:", os.Args[0], "AnchorPath PathWithinJail")
		os.Exit(1)
	}
	f, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening (%s)", err)
		os.Exit(1)
	}
	x, err := circuit.TryDial(f.Owner(), "acid")
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

	r := x.Call("JailOpen", os.Args[2])
	if r[1] != nil {
		fmt.Fprintf(os.Stderr, "Open problem: %#v\n", r[1])
		os.Exit(1)
	}
	io.Copy(os.Stdout, file.NewFileClient(r[0].(circuit.X)))
}
