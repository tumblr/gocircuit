package transport

import (
	"circuit/use/circuit"
	"encoding/gob"
	"net"
	"sync"
)

// Addr maintains a single unique instance for each addr.
// Addr object uniqueness is required by the circuit.Addr interface.
type Addr struct {
	ID   circuit.WorkerID
	PID  int
	Addr *net.TCPAddr
}

func init() {
	gob.Register(&Addr{})
}

func NewAddr(id circuit.WorkerID, pid int, hostport string) (circuit.Addr, error) {
	a, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil {
		return nil, err
	}
	return &Addr{ID: id, PID: pid, Addr: a}, nil
}

func (a *Addr) Host() string {
	return a.Addr.IP.String()
}

func (a *Addr) String() string {
	return a.ID.String() + "@" + a.Addr.String()
}

func (a *Addr) WorkerID() circuit.WorkerID {
	return a.ID
}

type addrTabl struct {
	lk   sync.Mutex
	tabl map[circuit.WorkerID]*Addr
}

func makeAddrTabl() *addrTabl {
	return &addrTabl{tabl: make(map[circuit.WorkerID]*Addr)}
}

func (t *addrTabl) Normalize(addr *Addr) *Addr {
	t.lk.Lock()
	defer t.lk.Unlock()

	a, ok := t.tabl[addr.ID]
	if ok {
		return a
	}
	t.tabl[addr.ID] = addr
	return addr
}
