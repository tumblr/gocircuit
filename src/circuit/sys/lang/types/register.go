package types

import (
	"strconv"
)


var (
	ValueTabl *TypeTabl = makeTypeTabl() // Type table for values
	FuncTabl  *TypeTabl = makeTypeTabl() // Type table for functions
)

// RegisterValue registers the type of x with the type table.
// Types need to be registered before values can be imported.
func RegisterValue(value interface{}) {
	ValueTabl.Add(makeType(value))
}

// RegisterFunc ...
func RegisterFunc(fn interface{}) {
	t := makeType(fn)
	if len(t.Func) != 1 {
		panic("fn type must have exactly one method: " + strconv.Itoa(len(t.Func)))
	}
	FuncTabl.Add(t)
}
