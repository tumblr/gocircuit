package config

import "circuit/use/circuit"

// SparkConfig ...
type SparkConfig struct {
	ID       circuit.RuntimeID
	BindAddr string
	Host     string
	Anchor   []string
}

// DefaultSpark is the default configuration used for workers started from the command line, which
// are often not intended to be contacted back from other workers
var DefaultSpark = &SparkConfig{
	ID:       circuit.ChooseRuntimeID(),	// Pick a random worker ID
	BindAddr: "",				// Don't accept incoming circuit calls from other workers
	Host:     "",				// "
	Anchor:   []string{},			// Don't register within the anchor file system
}
