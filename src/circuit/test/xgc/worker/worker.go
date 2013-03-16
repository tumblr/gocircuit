// Package worker defines the worker function for this application
package worker

import (
	"circuit/use/circuit"
)

type Start struct{}

func (Start) Main(dummy circuit.X) {
}

func init() { 
	circuit.RegisterFunc(Start{}) 
}

type Dummy struct{}

func init() { circuit.RegisterValue(&Dummy{}) }

func (*Dummy) Ping() {}
