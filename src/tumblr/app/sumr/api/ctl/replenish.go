package ctl

import (
	"path"
	"strconv"
	"sync"
	"tumblr/circuit/use/circuit"
	"tumblr/app/sumr/api"
	"tumblr/circuit/kit/sched/limiter"
	"tumblr/circuit/use/anchorfs"
)

type Result struct {
	Config      *WorkerConfig
	Replenished bool
	Err         error
}

// dfile is the name of the durable file describing the SUMR server cluster
func Replenish(dfile string, c *Config) []*Result {
	var (
		lk     sync.Mutex
		lmtr   limiter.Limiter
	)
	r := make([]*Result, len(c.Workers))
	lmtr.Init(20)
	for i_, wcfg_ := range c.Workers {
		i, wcfg := i_, wcfg_
		lmtr.Go(
			func() { 
				re, err := replenishWorker(dfile, c, i) 
				lk.Lock()
				defer lk.Unlock()
				r[i] = &Result{Config: wcfg, Replenished: re, Err: err}
			},
		)
	}
	lmtr.Wait()
	return r
}

func replenishWorker(dfile string, c *Config, i int) (replenished bool, err error) {

	// Check if worker already running
	anchor := path.Join(c.Anchor, strconv.Itoa(i))
	dir, e := anchorfs.OpenDir(anchor)
	if e != nil {
		return false, e
	}
	_, files, err := dir.Files()
	if e != nil {
		return false, e
	}
	if len(files) > 0 {
		return false, nil
	}

	// If not, start a new worker
	retrn, _, err := circuit.Spawn(c.Workers[i].Host, []string{anchor}, dfile, c.Workers[i].Port, c.ReadOnly)
	if err != nil {
		return false, err
	}
	if retrn[1] != nil {
		err = retrn[1].(error)
		return false, err
	}

	return true, nil
}

 
// Function for starting an API instance
func (Start) Start(dfile string, port int, readOnly bool) (circuit.XPerm, error) {
	a, err := api.New(dfile, port, readOnly)
	if err != nil {
		return nil, err
	}
	return circuit.PermRef(a), nil
}

type Start struct{}
func init() {
	circuit.RegisterFunc(Start{})
}
