// 4dls lists the contents of the durable file system
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"circuit/use/durablefs"
	_ "circuit/load"
	"sort"
	"strings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		println("Usage:", os.Args[0], " DurablePathQuery")
		println("	Examples of DurablePathQuery: /dir, /dir/...")
		os.Exit(1)
	}
	var recurse bool
	q := strings.TrimSpace(flag.Args()[0])
	if strings.HasSuffix(q, "...") {
		q = q[:len(q)-len("...")]
		recurse = true
	}
	ls(q, recurse, false)
}

func ls(query string, recurse, short bool) {
	dir := durablefs.OpenDir(query)

	// Read dirs
	chldn := dir.Children()

	var entries Entries
	for _, info := range chldn {
		entries = append(entries, info)
	}
	sort.Sort(entries)

	// Print sub-directories
	for _, e := range entries {
		hasBody, hasChildren := ' ', ' '
		if e.HasBody {
			hasBody = '*'
		}
		if e.HasChildren {
			hasChildren = '/'
		}
		fmt.Printf("%c %s%c\n", hasBody, path.Join(query, e.Name), hasChildren)
		if recurse {
			ls(path.Join(query, e.Name), recurse, short)
		}
	}
}

type Entries []durablefs.Info

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Less(i, j int) bool {
	return e[i].Name < e[j].Name
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
