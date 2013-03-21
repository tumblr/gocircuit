package tcp

import "errors"

var (
	ErrClosed       = errors.New("already closed")
	ErrCollision    = errors.New("conn id collision")
	ErrNotSupported = errors.New("not supported")
	ErrAuth         = errors.New("authentication")
	ErrProto        = errors.New("protocol")
)
