package transport

import "circuit/use/circuit"

var (
	ErrEnd           = circuit.NewError("end")
	ErrAlreadyClosed = circuit.NewError("already closed")
	errCollision     = circuit.NewError("conn id collision")
	ErrNotSupported  = circuit.NewError("not supported")
	ErrAuth          = circuit.NewError("authentication")
)
