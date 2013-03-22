// 4dcat prints the contents of a file from the durable file system
package main

import (
	_ "circuit/load"
	"circuit/use/durablefs"
	"flag"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], " DurablePath")
		os.Exit(1)
	}
	_, err := durablefs.OpenFile(flag.Arg(0))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// XXX: Need to be able to read unregistered types
}
