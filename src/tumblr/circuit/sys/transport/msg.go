package transport

import (
	"encoding/gob"
	"tumblr/circuit/use/lang"
)

func init() {
	gob.Register(&welcomeMsg{})
	gob.Register(&openMsg{})
	gob.Register(&connMsg{})
	gob.Register(&linkMsg{})
}

// linkMsg is the link-level message format between to endpoints.
// The link level is responsible for ensuring reliable and ordered delivery in
// the presence of network partitions and lost connections, assuming an
// eventual successful reconnect.
type linkMsg struct {
	SeqNo   int64 // OPT: Use circular integer comparison and fewer bits
	AckNo   int64
	Payload interface{}
}

type welcomeMsg struct {
	ID  lang.RuntimeID   // Runtime ID of sender
	PID int              // Process ID of sender runtime
}

type openMsg struct {
	ID connID
}

type connMsg struct {
	ID      connID
	Payload interface{}
}
