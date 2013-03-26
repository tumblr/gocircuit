package anchorfs

import (
	"circuit/kit/join"
	"circuit/use/circuit"
)

var link = join.SetThenGet{Name: "anchor file system"}

// Bind is used internally to bind an implementation of this package to the public methods of this package
func Bind(v interface{}) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

// CreateFile creates a new ephemeral file in the anchor directory anchor and saves the worker address addr in it.
func CreateFile(anchor string, addr circuit.Addr) error {
	return get().CreateFile(anchor, addr)
}

// Created returns a slive of anchor directories within which this worker has created files with CreateFile.
func Created() []string {
	return get().Created()
}

// OpenDir opens the anchor directory anchor
func OpenDir(anchor string) (Dir, error) {
	return get().OpenDir(anchor)
}

// OpenFile opens the anchor file anchor
func OpenFile(anchor string) (File, error) {
	return get().OpenFile(anchor)
}
