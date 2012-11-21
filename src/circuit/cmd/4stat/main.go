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
		println("Usage:", os.Args[0], "AnchorPath ServiceName")
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

	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "Worker disappeared during call (%#v)", p)
			os.Exit(1)
		}
	}()

	retrn := x.Call("StatServiceOnBehalf", os.Args[2])
	fmt.Println(retrn[0].(string))
}

type Stringer interface {
	String() string
}
