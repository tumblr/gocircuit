package client

import (
	"log"
	"math"
	"sync"
	"time"
	"tumblr/app/sumr"
	"tumblr/app/sumr/server/ctl"
	"circuit/use/circuit"
	"circuit/kit/sched/limiter"
	"tumblr/struct/xor"
)

// TODO: Enforce read only

// Client ..
type Client struct {
	dfile      string
	readOnly   bool
	checkpoint *ctl.Checkpoint
	lmtr       limiter.Limiter	// Global client rate limiter
	lk         sync.Mutex
	metric     xor.Metric		// Items in the metric are shard
}

type shard struct {
	Key    sumr.Key
	Server circuit.XPerm
}

// ID implements xor.Metric.Point
func (s *shard) ID() xor.ID {
	return xor.ID(s.Key)
}

// dfile is the filename of the node in Durable FS where the service keeps its
// checkpoint structure.
func New(dfile string, readOnly bool) (*Client, error) {
	cli := &Client{dfile: dfile, readOnly: readOnly}
	cli.lmtr.Init(50)

	var err error
	if cli.checkpoint, err = ctl.ReadCheckpoint(dfile); err != nil {
		return nil, err
	}

	// Compute metric space
	for _, x := range cli.checkpoint.Workers {
		cli.addServer(x)
	}
	return cli, nil
}

func (cli *Client) addServer(x *ctl.WorkerCheckpoint) {
	cli.lk.Lock()
	defer cli.lk.Unlock()
	cli.metric.Add(&shard{x.Key, x.Server})
}

func (cli *Client) Add(updateTime time.Time, key sumr.Key, value float64) (result float64) {

	// Per-client rate-limiting
	cli.lmtr.Open()
	defer cli.lmtr.Close()

	// Recover from dead shard panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("dead shard: %s", err)
			// XXX: Take more comprehensive action here
			result = math.NaN()
		}
	}()

	cli.lk.Lock()
	server := cli.metric.Nearest(xor.ID(key), 1)[0].(*ctl.WorkerCheckpoint).Server
	cli.lk.Unlock()

	retrn := server.Call("Add", updateTime, key, value)
	return retrn[0].(float64)
}

// AddRequest captures the input parameters for a Sumr ADD request
type AddRequest struct {
	UpdateTime time.Time
	Key        sumr.Key
	Value      float64
}

func (cli *Client) AddBatch(batch []AddRequest) []float64 {
	var lk sync.Mutex
	r := make([]float64, len(batch))

	blmtr := limiter.New(10)
	for i_, a_ := range batch {
		i, a := i_, a_
		blmtr.Go(func() {
			q := cli.Add(a.UpdateTime, a.Key, a.Value)
			lk.Lock()
			r[i] = q
			lk.Unlock()
		})
	}
	blmtr.Wait()
	return r
}

func (cli *Client) Sum(key sumr.Key) (result float64) {
	cli.lmtr.Open()
	defer cli.lmtr.Close()

	// Recover from dead shard panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("dead shard: %s", err)
			result = math.NaN()
		}
	}()

	cli.lk.Lock()
	server := cli.metric.Nearest(xor.ID(key), 1)[0].(*ctl.WorkerCheckpoint).Server
	cli.lk.Unlock()

	retrn := server.Call("Sum", key)
	return retrn[0].(float64)
}

// SumRequest captures the input parameters for a Sumr ADD request
type SumRequest struct {
	Key        sumr.Key
}

func (cli *Client) SumBatch(batch []SumRequest) []float64 {
	var lk sync.Mutex
	r := make([]float64, len(batch))

	blmtr := limiter.New(10)
	for i_, a_ := range batch {
		i, a := i_, a_
		blmtr.Go(func() {
			q := cli.Sum(a.Key)
			lk.Lock()
			r[i] = q
			lk.Unlock()
		})
	}
	blmtr.Wait()
	return r
}
