package main

import (
	"flag"
	"os"
	"tumblr/circuit/boot"
)

var flagShow = flag.Bool("v", false, "Verbose mode")

func main() {
	flag.Parse()
	c := boot.Build
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
