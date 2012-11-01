package main

import (
	"fmt"
	"log"
	"os"
	"tumblr/circuit/use/anchorfs"
	_ "tumblr/circuit/load"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "AnchorFile")
		os.Exit(1)
	}
	file, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", file.Owner())
}
