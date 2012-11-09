package ctl

import (
	"encoding/gob"
	"time"
	"tumblr/struct/xor"
	"tumblr/app/sumr"
	"circuit/use/circuit"
	"circuit/use/durablefs"
)

func init() {
	gob.RegisterName("sumr.server.ctl.Config", &Config{})
	gob.RegisterName("sumr.server.ctl.WorkerConfig", &WorkerConfig{})
	gob.RegisterName("sumr.server.ctl.Checkpoint", &Checkpoint{})
	gob.RegisterName("sumr.server.ctl.WorkerCheckpoint", &WorkerCheckpoint{})
}

// Config
type Config struct {
	Anchor  string		// Anchor for the SUMR shard workers
	Workers []*WorkerConfig
}

type WorkerConfig struct {
	Host     circuit.Host
	DiskPath string
	Forget   time.Duration
}

// Checkpoint
type Checkpoint struct {
	Config  *Config
	Workers []*WorkerCheckpoint
}

type WorkerCheckpoint struct {
	Key     sumr.Key	// Keys are allocated dynamically, stored only here
	Runtime circuit.Addr
	Server  circuit.XPerm
	Host    circuit.Host
}

// ReadCheckpoint reads a checkpoint structure from the durable file dfile.
func ReadCheckpoint(dfile string) (*Checkpoint, error) {
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
	chk, ok := chk_[0].(*Checkpoint)
	if !ok {
		return nil, circuit.NewError("unexpected checkpoint value (%#v) of type (%T) in durable file %s", chk_[0], chk_[0], dfile)
	}
	return chk, nil
}

// Implements xor.Item
func (s *WorkerCheckpoint) ID() xor.ID {
	return xor.ID(s.Key)
}
