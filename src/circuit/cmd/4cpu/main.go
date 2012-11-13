// 4cpu causes a running worker to start CPU profiling for a specified interval, after which it writes the pprof file locally
package main

import (
	"fmt"
	"log"
	"os"
	_ "circuit/load"
	"circuit/use/anchorfs"
	"circuit/use/circuit"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage:", os.Args[0], "AnchorPath DurationSeconds")
		os.Exit(1)
	}
	
	// Parse duration
	dursec, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing duration (%s)\n", err)
		os.Exit(1)
	}
	dur := time.Duration(int64(dursec)*1e9)

	// Find anchor file
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

	// Connect to worker
	retrn := x.Call("CPUProfile", dur)
	if err, ok := retrn[1].(error); ok && err != nil {
		log.Printf("Problem obtaining CPU profile (%s)", err)
		os.Exit(1)
	}
	fmt.Println(string(retrn[0].([]byte)))
}
