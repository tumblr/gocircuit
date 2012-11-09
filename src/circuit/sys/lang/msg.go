package lang

import (
	"encoding/gob"
	"circuit/use/circuit"
	"circuit/sys/lang/types"
)

func init() {
	gob.Register(&exportedMsg{})
	// Func invokation-style commands
	gob.Register(&goMsg{})
	gob.Register(&callMsg{})
	gob.Register(&dialMsg{})
	gob.Register(&getPtrMsg{})
	gob.Register(&returnMsg{})
	// Value-passing internal commands
	gob.Register(&gotPtrMsg{})
	gob.Register(&dontReplyMsg{})
	gob.Register(&dropPtrMsg{})
	// Value-passing sub-messages
	gob.Register(&ptrMsg{})
	gob.Register(&ptrPtrMsg{})
	gob.Register(&permPtrMsg{})
	gob.Register(&permPtrPtrMsg{})
}

// Top-level messages

type exportedMsg struct {
	Value  []interface{}
	Stack string
}

// Execute a method call
type callMsg struct {
	ReceiverID handleID
	FuncID     types.FuncID
	In         []interface{}
}

// Fork a go routine
type goMsg struct {
	TypeID types.TypeID
	In     []interface{}
}

type returnMsg struct {
	Out []interface{}
	Err error
}

type getPtrMsg struct {
	ID handleID
}

type gotPtrMsg struct {
	ID handleID
}

// dontReplyMsg is dropped by the receiver and intentionally never replies to.
// It is used to sense the death of a runtime.
type dontReplyMsg struct{}

// dialMsg requests that the receiver send back a handle to its permanent.
type dialMsg struct{
	Service string
}

// The importer of a handle sends a release request to the exporter to
// notify them that the held object is no longer needed.
// This is part of the cross-runtime garbage collection mechanism.
type dropPtrMsg struct {
	ID handleID
}

// ptrMsg carries ...
type ptrMsg struct {
	ID      handleID
	TypeID  types.TypeID
}

func (*ptrMsg) IsX() {}

func (*ptrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (*ptrMsg) String() string {
	panic("not for call")
}

// ptrPtrMsg carries ...
type ptrPtrMsg struct {
	ID  handleID
	Src circuit.Addr
}

func (*ptrPtrMsg) IsX() {}

func (*ptrPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (*ptrPtrMsg) String() string {
	panic("not for call")
}

// permPtrMsg carries ...
type permPtrMsg struct {
	ID      handleID
	TypeID  types.TypeID
}

func (*permPtrMsg) IsX() {}

func (*permPtrMsg) IsXPerm() {}

func (*permPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (*permPtrMsg) String() string {
	panic("not for call")
}

// permPtrPtrMsg carries a serialized parmenent x-pointer from a sender to a receiver,
// where the value pointed to is not owned by the sender.
type permPtrPtrMsg struct {
	ID     handleID
	TypeID types.TypeID
	Src    circuit.Addr
}

func (*permPtrPtrMsg) IsX() {}

func (*permPtrPtrMsg) IsXPerm() {}

func (*permPtrPtrMsg) Call(proc string, in ...interface{}) []interface{} {
	panic("hack: not meant to be used")
}

func (*permPtrPtrMsg) String() string {
	panic("not for call")
}
