package main

import (
	"flag"
	"os"
	//"circuit/load/config"
	"fmt"
	"path"
)

func init() {
	flag.Usage = func() {
		_, prog := path.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [RuntimeID]\n", prog)
		fmt.Fprintf(os.Stderr,
`
4hardkill kills worker processes pertaining to the contextual circuit on all
hosts supplied on standard input and separated by new lines.

Instead of using the in-circuit facilities to do so, this utility logs directly
into the target hosts (using ssh), finds and kills relevant processes using
POSIX-level facilities.

If a RuntimeID is specified, only the worker having the ID in question is killed.
`)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	flag.Usage()
	panic("not implemented")
}
