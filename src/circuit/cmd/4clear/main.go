// 4clear deletes the jails of workers that are no longer alive, on a list of hosts specified one per-line on standard input
package main

import (
	"bufio"
	"circuit/kit/posix"
	"circuit/load/config"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func init() {
	flag.Usage = func() {
		_, prog := path.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s\n", prog)
		fmt.Fprintf(os.Stderr,
			`
4clear deletes the jails of workers that are no longer alive, 
on all hosts specified one per-line on standard input.
`)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Read target hosts from standard input
	var hosts []string
	buf := bufio.NewReader(os.Stdin)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			line = strings.TrimSpace(line)
			hosts = append(hosts, line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading target hosts (%s)", err)
			os.Exit(1)
		}
	}

	// Log into each host and kill pertinent workers, using POSIX kill
	for _, h := range hosts {
		println("Clearing dead worker jails on", h)
		clearSh := fmt.Sprintf("%s %s\n", config.Config.Install.ClearHelperPath(), config.Config.Install.JailDir())
		_, stderr, err := posix.Exec("ssh", "", clearSh, h, "sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem while clearing jails on %s (%s)\n", h, err)
			fmt.Fprintf(os.Stderr, "Remote clear-helper error output:\n%s\n", stderr)
		}
	}
}
