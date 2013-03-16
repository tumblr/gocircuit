// 4hardkill kills worker processes on a given host using out-of-band UNIX-level facilities
package main

import (
	"bufio"
	"flag"
	"os"
	"fmt"
	"path"
	"circuit/kit/posix"
	"circuit/use/circuit"
	"circuit/load/config"
	"io"
	"strings"
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

	// Parse RuntimeID argument
	var (
		err    error
		id     circuit.RuntimeID
		withID bool
	)
	if flag.NArg() == 1 {
		id, err = circuit.ParseRuntimeID(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem parsing runtime ID (%s)\n", err)
			os.Exit(1)
		}
		withID = true
	} else if flag.NArg() != 0 {
		flag.Usage()
	}

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
		println("Hard-killing circuit worker(s) on", h)
		var killSh string
		if withID {
			killSh = fmt.Sprintf("ps ax | grep -i %s | grep -v grep | awk '{print $1}' | xargs kill -KILL\n", id.String())
		} else {
			killSh = fmt.Sprintf("ps ax | grep -i %s | grep -v grep | awk '{print $1}' | xargs kill -KILL\n", config.Config.Install.Binary)
		}
		_, stderr, err := posix.Exec("ssh", "", killSh, h, "sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem while killing workers on %s (%s)\n", h, err)
			fmt.Fprintf(os.Stderr, "Remote shell error output:\n%s\n", stderr)
		}
	}
}
