package main

import (
	"flag"
	"os"
	"tumblr/TUMBLR/config"
)

var flagShow = flag.Bool("v", false, "Verbose mode")

func main() {
	flag.Parse()
	c := config.Build
	println("Building circuit on", c.Host)
	c.Show = *flagShow
	if err := Build(c); err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println("Done.")
}
