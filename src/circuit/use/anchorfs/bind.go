package anchorfs

import (
	"circuit/use/circuit"
	"circuit/kit/join"
)

var link = join.SetThenGet{Name: "anchor file system"}

func Bind(v fs) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

func CreateFile(anchor string, addr circuit.Addr) error {
	return get().CreateFile(anchor, addr)
}

func Created() []string {
	return get().Created()
}

func OpenDir(anchor string) (Dir, error) {
	return get().OpenDir(anchor)
}

func OpenFile(anchor string) (File, error) {
	return get().OpenFile(anchor)
}
