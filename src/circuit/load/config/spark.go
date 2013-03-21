package config

import "circuit/use/circuit"

// SparkConfig ...
type SparkConfig struct {
	ID       circuit.WorkerID
	BindAddr string
	Host     string
	Anchor   []string
}

// DefaultSpark is the default configuration used for workers started from the command line, which
// are often not intended to be contacted back from other workers
var DefaultSpark = &SparkConfig{
	ID:       circuit.ChooseWorkerID(),
	BindAddr: "",         // Don't accept incoming circuit calls from other workers
	Host:     "",         // "
	Anchor:   []string{}, // Don't register within the anchor file system
}
