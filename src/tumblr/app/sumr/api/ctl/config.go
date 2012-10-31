package ctl

import (
	"encoding/gob"
	"tumblr/circuit"
)

func init() {
	// Automate this custom-named registration process in tumblr/circuit
	gob.RegisterName("sumr.api.ctl.Config", &Config{})
	gob.RegisterName("sumr.api.ctl.WorkerConfig", &WorkerConfig{})
}

// Config
type Config struct {
	Anchor   string		// Anchor for the SUMR api workers
	ReadOnly bool
	Workers  []*WorkerConfig
}

type WorkerConfig struct {
	Host circuit.Host
	Port int
}
