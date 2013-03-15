// 4crossbuild automates the process of cross-building a circuit application remotely
package main

import (
	"flag"
	"os"
	"circuit/load/config"
)

var flagShow = flag.Bool("v", false, "Verbose mode")

func main() {
	flag.Parse()
	c := config.Config.Build
	c.Binary = config.Config.Install.Binary
	if c == nil {
		println("Circuit build configuration not specified in environment")
		os.Exit(1)
	}
	println("Building circuit on", c.Host)
	c.Show = *flagShow
	if err := Build(c); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println("Done.")
}
