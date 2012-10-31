package zanchorfs

import (
	"tumblr/circuit/use/lang"
)

type File struct {
	owner lang.Addr
}

func (f *File) Owner() lang.Addr {
	return f.owner
}
