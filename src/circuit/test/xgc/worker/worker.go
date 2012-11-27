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
