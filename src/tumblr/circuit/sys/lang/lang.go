package lang

import "tumblr/circuit/use/circuit"

// _ref wraps a user object, indicating to the runtime that the user has
// elected to send this object as a ptr across runtimes.
type _ref struct {
	value interface{}
}

func (*_ref) String() string {
	return "XREF"
}

func (*_ref) IsX() {}

func (*_ref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on ref")
}

type _permref struct {
	value interface{}
}

func (*_permref) String() string {
	return "XPERMREF"
}

func (*_permref) IsX() {}

func (*_permref) IsXPerm() {}

func (*_permref) Call(proc string, in ...interface{}) []interface{} {
	panic("call on permref")
}

// Ref annotates a user value v, so that if the returned value is consequently
// passed cross-runtime, the runtime will pass v as via a cross-runtime pointer
// rather than by value.
func (*Runtime) Ref(v interface{}) circuit.X {
	return Ref(v)
}

func Ref(v interface{}) circuit.X {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		return v
	case *_ref:
		return v
	case *_permptr:
		return v
	case *_permref:
		panic("applying ref on permref")	
	}
	return &_ref{v}
}

func (*Runtime) PermRef(v interface{}) circuit.XPerm {
	return PermRef(v)
}

func PermRef(v interface{}) circuit.XPerm {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case *_ptr:
		panic("permref on ptr")
	case *_ref:
		panic("permref on ref")
	case *_permptr:
		return v
	case *_permref:
		return v
	}
	return &_permref{v}
}
