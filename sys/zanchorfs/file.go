package zanchorfs

import (
	"circuit/use/circuit"
)

type File struct {
	owner circuit.Addr
}

func (f *File) Owner() circuit.Addr {
	return f.owner
}
