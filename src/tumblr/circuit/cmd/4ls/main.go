package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"tumblr/circuit/use/anchorfs"
	_ "tumblr/TUMBLR/load"
)

var flagFull = flag.Bool("f", true, "Print out full path of files and directories")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], "[-f] AnchorPath")
		os.Exit(1)
	}
	query := flag.Args()[0]
	dir, err := anchorfs.OpenDir(query)
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}
	dirs, err := dir.Dirs()
	if err != nil {
		log.Printf("Problem listing directories (%s)", err)
		os.Exit(1)
	}
	_, files, err := dir.Files()
	if err != nil {
		log.Printf("Problem listing files (%s)", err)
		os.Exit(1)
	}
	// Print sub-directories
	for _, d := range dirs {
		if *flagFull {
			fmt.Println(path.Join(query, d))
		} else {
			fmt.Printf("/%s\n", d)
		}
	}
	// Print files
	for f, _ := range files {
		if *flagFull {
			fmt.Println(path.Join(query, f.String()))
		} else {
			fmt.Printf("%s\n", f)
		}
	}
}
