package main

import (
	"fmt"
	"log"
	"os"
	_ "tumblr/circuit/boot"
	"tumblr/circuit/use/anchorfs"
	"tumblr/circuit/use/circuit"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "AnchorPath")
		os.Exit(1)
	}
	file, err := anchorfs.OpenFile(os.Args[1])
	if err != nil {
		log.Printf("Problem opening (%s)", err)
		os.Exit(1)
	}
	x, err := circuit.TryDial(file.Owner(), "acid")
	if err != nil {
		log.Printf("Problem dialing acid service (%s)", err)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			log.Printf("Worker disappeared during call (%#v)", p)
			os.Exit(1)
		}
	}()

	retrn := x.Call("RuntimeProfile", "goroutine", 1)
	if err, ok := retrn[1].(error); ok && err != nil {
		log.Printf("Problem obtaining runtime profile (%s)", err)
		os.Exit(1)
	}
	fmt.Println(string(retrn[0].([]byte)))
}
