// 4top displays real-time vitals (cpu, mem, io) of circuit deployments at various anchor granularities (file, directory, subtree)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	_ "circuit/load"
	"sort"
	"strings"
)

var flagShort = flag.Bool("s", false, "Do not print full path")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], "AnchorPath")
		println("	Examples of AnchorPath: /host, /host/...")
		println(
`4top displays real-time vitals (cpu, mem, io) of circuit deployments at
various anchor granularities (file, directory, subtree).
`)
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	top(q, recurse, *flagShort)
}

func fileMapToSlice(m map[circuit.RuntimeID]anchorfs.File) []string {
	var r []string
	for id, _ := range m {
		r = append(r, id.String())
	}
	return r
}

type info struct {
	Worker string
	CPU    float64
	Mem    int64
}

func top(query string, recurse, short bool) []*info {
	dir, err := anchorfs.OpenDir(query)
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}

	// Read dirs
	dirs, err := dir.Dirs()
	if err != nil {
		log.Printf("Problem listing directories (%s)", err)
		os.Exit(1)
	}
	sort.Strings(dirs)

	// Read files
	_, filesMap, err := dir.Files()
	if err != nil {
		log.Printf("Problem listing files (%s)", err)
		os.Exit(1)
	}
	files := fileMapToSlice(filesMap)
	sort.Strings(files)

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
	for _, f := range files {
		if !*flagShort {
			fmt.Println(path.Join(query, f))
		} else {
			fmt.Printf("%s\n", f)
		}
	}
}
