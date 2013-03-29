// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 4deploy installs locally-available circuit binaries to a cluster of hosts supplied to standard input, one host per line
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
	Install(config.Config.Deploy, config.Config.Build, hosts)
	println("Done.")
}
