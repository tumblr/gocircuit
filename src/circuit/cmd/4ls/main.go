package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"circuit/use/anchorfs"
	_ "circuit/load"
	"strings"
)

var flagShort = flag.Bool("s", false, "Do not print full path")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], "[-s] AnchorPath")
		println("	-s Do not print full path")
		println("	Examples of AnchorPath: /host, /host/...")
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	ls(q, recurse, *flagShort)
}

func ls(query string, recurse, short bool) {
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
		if !*flagShort {
			fmt.Println(path.Join(query, d))
		} else {
			fmt.Printf("/%s\n", d)
		}
		if recurse {
			ls(path.Join(query, d), recurse, short)
		}
	}
	// Print files
	for f, _ := range files {
		if !*flagShort {
			fmt.Println(path.Join(query, f.String()))
		} else {
			fmt.Printf("%s\n", f)
		}
	}
}
