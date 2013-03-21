// 4jailtail redirects the output of running 'tail -f' on a file inside a worker's local file system jail
package main

import (
	teleio "circuit/kit/tele/io"
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"fmt"
	"io"
	"os"
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

	tailViaSSH(f.Owner(), os.Args[2])
}

func tailViaCircuit(addr circuit.Addr, jailpath string) {

	x, err := circuit.TryDial(addr, "acid")
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

	r := x.Call("JailTail", jailpath)
	if r[1] != nil {
		fmt.Fprintf(os.Stderr, "Open problem: %s\n", r[1].(error))
		os.Exit(1)
	}

	io.Copy(os.Stdout, teleio.NewClient(r[0].(circuit.X)))

	/*
		tailr := teleio.NewClient(r[0].(circuit.X))
		for {
			p := make([]byte, 1e3)
			n, err := tailr.Read(p)
			if err != nil {
				println(err.Error(),"+++")
				break
			}
			println("n=", n)
		}*/
}
