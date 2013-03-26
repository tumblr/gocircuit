package config

import "circuit/use/circuit"

// SparkConfig captures a few worker startup parameters that can be configured on each execution
type SparkConfig struct {
	// ID is the ID of the worker instance
	ID       circuit.WorkerID

	// BindAddr is the network address the worker will listen to for incoming connections
	BindAddr string

	// Host is the host name of the hosting machine
	Host     string

	// Anchor is the set of anchor directories that the worker registers with
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
