// Package durablefs exposes the programming interface to a global file system for storing cross-values
package durablefs

import (
	"time"
	"circuit/kit/join"
	"circuit/use/circuit"
)

var link = join.SetThenGet{Name: "durable file system"}

func Bind(v fs) {
	link.Set(v)
}

func get() fs {
	return link.Get().(fs)
}

func OpenFile(name string) (File, error) {
	return get().OpenFile(name)
}

func CreateFile(name string) (File, error) {
	return get().CreateFile(name)
}

func Remove(name string) error {
	return get().Remove(name)
}

func OpenDir(name string) Dir {
	return get().OpenDir(name)
}

// The fs, Dir and Conn interfaces return an error when an error reflects an
// expected user-level condition. E.g. CreateFile will return an error if the
// file exists. This is often a valid execution path. On the other hand, fs
// operations panic if a user-independent condition, like a network outage,
// occurs.

type fs interface {

	// File operations
	OpenFile(string) (File, error)
	CreateFile(string) (File, error)

	// File or directory
	Remove(string) error

	// Dir operations
	OpenDir(string) Dir
}

type Dir interface {
	Path() string
	Children() (children map[string]struct{})
	Change() (children map[string]struct{})
	Expire(expire time.Duration) (children map[string]struct{})
	Close()
}

type File interface {
	Read() ([]interface{}, error)
	Write(...interface{}) error
	Close() error
}

var ErrParse = circuit.NewError("parse")
