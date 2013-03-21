package api

import "encoding/gob"

func init() {
	gob.Register(&Config{})
	gob.Register(&WorkerConfig{})
}

// Config specifies a cluster of HTTP API servers
type Config struct {
	Anchor   string          // Anchor for the sumr API workers
	ReadOnly bool            // Reject requests resulting in change
	Workers  []*WorkerConfig // Specification of service workers
}

// WorkerConfig specifies an individual API server
type WorkerConfig struct {
	Host string // Host is the circuit hostname where the worker is to be deployed
	Port int    // Port is the port number when the HTTP API server is to listen
}
