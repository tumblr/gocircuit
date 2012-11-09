package lang

import (
	"log"
	"sync"
	"tumblr/circuit/use/circuit"
	"tumblr/circuit/sys/lang/prof"
	"tumblr/circuit/sys/lang/types"
)

// Runtime represents that state of the circuit program at the present moment.
// This state can change in two ways: by a 'linguistic' action ...
type Runtime struct {
	dialer  circuit.Transport
	exp     *expTabl
	imp     *impTabl
	srv     srvTabl
	blk     sync.Mutex
	boot    interface{}
	lk      sync.Mutex
	live    map[circuit.Addr]struct{}  // Set of peers we monitor for liveness
	prof    *prof.Profile

	dlk     sync.Mutex
	dallow  bool
	daemon  bool
}

func New(t circuit.Transport) *Runtime {
	r := &Runtime{
		dialer: t,
		exp:    makeExpTabl(types.ValueTabl),
		imp:    makeImpTabl(types.ValueTabl),
		live:   make(map[circuit.Addr]struct{}),
		prof:   prof.New(),
	}
	r.srv.Init()
	go func() {
		for {
			r.accept(t)
		}
	}()
	r.Listen("acid", &acid{})
	return r
}

func (r *Runtime) XAddr() circuit.Addr {
	return r.dialer.Addr()
}

func (r *Runtime) SetBoot(v interface{}) {
	r.blk.Lock()
	defer r.blk.Unlock()
	if v != nil {
		types.RegisterType(v)
	}
	r.boot = v
}

func (r *Runtime) accept(l circuit.Listener) {
	conn := l.Accept()
	req, err := conn.Read()
	if err != nil {
		panic(err)
		return
	}
	// Demux the request
	go func() {
		// Importing reptr variables involves waiting on other runtimes,
		// we fork request handling to dedicated go routines.
		// No rate-limiting/throttling is performed in the circuit.
		// It is the responsibility of Listener and/or the user app logic to
		// keep the runtime from contending.
		switch q := req.(type) {
		case *goMsg:
			r.serveGo(q, conn)
		case *dialMsg:
			r.serveDial(q, conn)
		case *callMsg:
			r.serveCall(q, conn)
		case *dropPtrMsg:
			r.serveDropPtr(q, conn)
		case *getPtrMsg:
			r.serveGetPtr(q, conn)
		case *dontReplyMsg:
			// Don't reply. Intentionally don't close the conn.
			// It will close when the process dies.
		default:
			log.Printf("unknown request %v", req)
		}
	}()
}
