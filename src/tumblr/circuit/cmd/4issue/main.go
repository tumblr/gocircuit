package main

import (
	"fmt"
	"os"
	"tumblr/circuit/use/issuefs"
	_ "tumblr/TUMBLR/load"
)

func usage() {
	println("Usage:", os.Args[0], "(ls | resolve ID | subscribe Email | unsubscribe Email)")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "subscribe":
		if len(os.Args) != 3 {
			usage()
		}
		if err := issuefs.Subscribe(os.Args[2]); err != nil {
			println("Email already subscribed")
			os.Exit(1)
		}
	case "unsubscribe":
		if len(os.Args) != 3 {
			usage()
		}
		if err := issuefs.Unsubscribe(os.Args[2]); err != nil {
			println("Email not subscribed")
			os.Exit(1)
		}
	case "ls":
		issues := issuefs.List()
		for _, i := range issues {
			fmt.Printf("%s\n", i.String())
		}
	case "resolve":
		if len(os.Args) != 3 {
			usage()
		}
		id, err := issuefs.ParseID(os.Args[2])
		if err != nil {
			println("Issue ID did not parse correctly")
			os.Exit(1)
		}
		if err = issuefs.Resolve(id); err != nil {
			println("No issue with this id")
			os.Exit(1)
		}
	default:
		usage()
	}
}
