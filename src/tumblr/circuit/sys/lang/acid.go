package lang

import (
	"bytes"
	"log"
	"runtime/pprof"
	"time"
	"tumblr/circuit/use/lang"
)

type acid struct{}

/*
func (s *acid) Stat(runtime.Frame) *profile.WorkerStat {
	return s.profile.Stat()
}
*/

// Ping is a nop. Its intended use is as a basic check whether a worker is still alive.
func (s *acid) Ping() {}

// RuntimeProfile exposes the Go runtime profiling framework of this worker
func (s *acid) RuntimeProfile(name string, debug int) ([]byte, error) {
	prof := pprof.Lookup(name)
	if prof == nil {
		return nil, lang.NewError("no such profile")
	}
	var w bytes.Buffer
	if err := prof.WriteTo(&w, debug); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *acid) CPUProfile(duration time.Duration) ([]byte, error) {
	if duration > time.Hour {
		return nil, lang.NewError("cpu profile duration exceeds 1 hour")
	}
	var w bytes.Buffer
	if err := pprof.StartCPUProfile(&w); err != nil {
		return nil, err
	}
	log.Printf("cpu profiling for %d sec", duration / 1e9)
	time.Sleep(duration)
	pprof.StopCPUProfile()
	return w.Bytes(), nil
}
