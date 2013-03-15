// 4clear-helper is used internally by 4clear
package main

import (
	"os"
	"circuit/kit/lockfile"
	"circuit/use/circuit"
	"fmt"
	"path"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "JailDir")
		os.Exit(1)
	}
	jailDir := os.Args[1]

	jail, err := os.Open(jailDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open jail directory (%s)\n", err)
		os.Exit(1)
	}
	defer jail.Close()
	fifi, err := jail.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read jail directory (%s)\n", err)
		os.Exit(1)
	}

	for _, fi := range fifi {
		if !fi.IsDir() {
			continue
		}
		if _, err := circuit.ParseRuntimeID(fi.Name()); err != nil {
			continue
		}

		workerJail := path.Join(jailDir, fi.Name())
		println("Clearing", workerJail)
		l, err := lockfile.Create(path.Join(workerJail, "lock"))
		if err != nil {
			// This worker is alive; still holding lock; move on
			println(err.Error())
			continue
		}
		l.Release()
		if err := os.RemoveAll(workerJail); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot remove worker jail %s (%s)\n", workerJail, err)
			os.Exit(1)
		}
	}
}
