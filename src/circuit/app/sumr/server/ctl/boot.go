package ctl

import (
	"circuit/app/sumr"
	"circuit/app/sumr/server"
	"circuit/kit/sched/limiter"
	"circuit/use/circuit"
	"log"
	"sync"
	"tumblr/struct/xor"
)

// Boot launches a SUMR instance as specified in c.
// In its lifetime (across failures and restarts), a service is only booted
// once.  Reviving service shards in response to external events is done via
// the various control functions. The services' current persistent state is
// stored in the durable FS under the name dfile. Future maintenance to the
// service is possible due to this durable state.
//
func Boot(c *Config) *Checkpoint {
	s := &Checkpoint{Config: c}
	s.Workers = boot(c.Anchor, c.Workers)
	return s
}

// Boot starts a SUMR shard server on each host specified in cluster, and returns
// a list of shards and respective keys and a corresponding list of runtime processes.
//
func boot(anchor string, shard []*WorkerConfig) []*WorkerCheckpoint {
	var (
		lk     sync.Mutex
		lmtr   limiter.Limiter
		shv    []*WorkerCheckpoint
		metric xor.Metric // Used to allocate initial keys in a balanced fashion
	)
	lmtr.Init(20)
	shv = make([]*WorkerCheckpoint, len(shard))
	for i_, sh_ := range shard {
		i, sh := i_, sh_
		xkey := metric.ChooseMinK(5)
		lmtr.Go(
			func() {
				x, addr, err := bootShard(anchor, sh)
				if err != nil {
					log.Printf("sumr shard boot on %s error (%s)", sh.Host, err)
					return
				}
				lk.Lock()
				defer lk.Unlock()
				shv[i] = &WorkerCheckpoint{
					Key:     sumr.Key(xkey),
					Runtime: addr,
					Server:  x,
					Host:    sh.Host,
				}
			},
		)
	}
	lmtr.Wait()
	return shv
}

func bootShard(anchor string, sh *WorkerConfig) (x circuit.XPerm, addr circuit.Addr, err error) {

	retrn, addr, err := circuit.Spawn(sh.Host, []string{anchor}, server.main{}, sh.DiskPath, sh.Forget)
	if retrn[1] != nil {
		err = retrn[1].(error)
		return nil, nil, err
	}

	return retrn[0].(circuit.XPerm), addr, nil
}
