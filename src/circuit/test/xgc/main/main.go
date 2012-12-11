package main

import (
	_ "circuit/load"
	"circuit/test/xgc/worker"
	"circuit/use/circuit"
	"circuit/use/n"
	"runtime"
	_ "circuit/kit/debug/ctrlc"
)

// TODO: Make sure finalizer called BECAUSE worker died or worker asked us to release handle

func main() {
	ch := make(chan int)
	d := &worker.Dummy{}
	runtime.SetFinalizer(d, func(h *worker.Dummy) {
		println("finalizing dummy")
		close(ch)
	})

	// Test: 
	//	Spawn a worker and pass an x-pointer to it; 
	//	Worker proceeds to die right away;
	//	Check that finalizer of local dummy called when local runtime notices remote is dead
	_, addr, err := circuit.Spawn(n.ParseHost("localhost"), []string{"/xgc"}, worker.Start{}, circuit.Ref(d))
	if err != nil {
		panic(err)
	}
	d = nil // Make sure we are not holding the object
	runtime.GC()

	println(addr.String())
	println("Waiting for finalizer call ...")
	<-ch
	println("Success")
}
