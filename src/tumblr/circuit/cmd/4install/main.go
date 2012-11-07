// 4install ships and installs a built circuit from the local shipping
// directory to a group of hosts specified on standard input in the form of one
// hostname (no port number) per line
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"tumblr/circuit/load/config"
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
