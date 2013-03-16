package server

import (
	"bytes"
	"circuit/app/sumr"
	"circuit/use/circuit"
	"circuit/use/durablefs"
	"encoding/gob"
	"fmt"
	"time"
	"tumblr/struct/xor"
)

func init() {
	gob.Register(&Config{})
	gob.Register(&WorkerConfig{})
	gob.Register(&checkpoint{})
	gob.Register(&workerCheckpoint{})
}

// Config specifies a cluster of sumr shard servers
type Config struct {
	Anchor  string          // Anchor directory for the sumr shard workers
	Workers []*WorkerConfig // List of workers
}

// WorkerConfig specifies a configuration for an individual sumr shard worker
type WorkerConfig struct {
	Host     string        // Host is the circuit hostname where the worker is to be deployed
	DiskPath string        // DiskPath is a local directory to be used for persisting the shard
	Forget   time.Duration // Key-value pairs older than Forget will be evicted from memory and unavailable for querying
}

// Checkpoint represents the runtime configuration of a live sumr database
type checkpoint struct {
	Config  *Config             // Config is the configuration used to start the database service
	Workers []*workerCheckpoint // Workers is a list of the runtime configuration of all shard workers
}

// String returns a textual representation of this checkpoint
func (s *checkpoint) String() string {
	var w bytes.Buffer
	for i, shc := range s.Config.Workers {
		srvstr := "•"
		key := "•"
		if shs := s.Workers[i]; shs != nil {
			srvstr = shs.Server.String()
			key = shs.Key.String()
		}
		fmt.Fprintf(&w, "KEY=%s SERVER=%s HOST=%s DISK=%s FORGET=%s\n", key, srvstr, shc.Host, shc.DiskPath, shc.Forget)
	}
	return string(w.Bytes())
}

// WorkerCheckpoint represents the runtime configuration of a shard worker
type workerCheckpoint struct {
	Key    sumr.Key      // Key is the key of the shard; keys are assigned dynamically after worker startup
	Addr   circuit.Addr  // Addr is the address of the live worker shard
	Server circuit.XPerm // Server is a permanent cross-interface to the shard receiver
	Host   string        // Host is the circuit hostname where this worker is executing
}

// ReadCheckpoint reads a checkpoint structure from the durable file dfile.
func readCheckpoint(dfile string) (*checkpoint, error) {
	// Fetch service info from durable fs
	f, err := durablefs.OpenFile(dfile)
	if err != nil {
		return nil, err
	}
	chk_, err := f.Read()
	if err != nil {
		return nil, err
	}
	if len(chk_) == 0 {
		return nil, circuit.NewError("no values in checkpoint durable file " + dfile)
	}
	chk, ok := chk_[0].(*checkpoint)
	if !ok {
		return nil, circuit.NewError("unexpected checkpoint value (%#v) of type (%T) in durable file %s", chk_[0], chk_[0], dfile)
	}
	return chk, nil
}

// ID returns the XOR-metric ID of the shard underlying this checkpoint
func (s *workerCheckpoint) ID() xor.Key {
	return xor.Key(s.Key)
}
