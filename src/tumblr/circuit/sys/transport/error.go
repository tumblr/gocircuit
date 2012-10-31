package transport

import (
	"tumblr/circuit/use/lang"
)

var (
	ErrEnd            = lang.NewError("end")
	ErrAlreadyClosed  = lang.NewError("already closed")
	errCollision      = lang.NewError("conn id collision")
	ErrNotSupported   = lang.NewError("not supported")
)
