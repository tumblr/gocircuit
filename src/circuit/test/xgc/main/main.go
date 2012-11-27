package main

import (
	_ "circuit/load"
	"circuit/test/xgc/worker"
	"circuit/use/circuit"
	"circuit/use/n"
	"runtime"
)

// TODO:
//	Spawn does not return
//	Make sure finalizer called BECAUSE worker died or worker asked us to release handle

type Dummy struct{}
func init() { circuit.RegisterValue(&Dummy{}) }

func (*Dummy) Ping() {}

func main() {
	ch := make(chan int)
	d := &Dummy{}
	runtime.SetFinalizer(d, func(h *Dummy) {
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

	println(addr.String())
	println("Waiting for finalizer call ...")
	<-ch
	println("Success")
}
