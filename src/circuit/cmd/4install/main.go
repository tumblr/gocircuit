// 4install installs locally-available circuit binaries to a cluster of hosts supplied to standard input, one host per line
package main

import (
	"bufio"
	"circuit/load/config"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
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
	println("Installing circuit.")
	Install(config.Config.Install, config.Config.Build, hosts)
	println("Done.")
}
