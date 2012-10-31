// Package fs defines a general interface for file-systems.
package fs

import (
	"errors"
	"os"
)

type FS interface {
	Open(name string) (File, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Create(name string) (File, error)
	Remove(name string) error
	Rename(oldname, newname string) error
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string) error
	MkdirAll(name string) error
}

type File interface {
	Close() error
	Stat() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Read([]byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
	Truncate(size int64) error
	Write([]byte) (int, error)
	Sync() error
}

var (
	ErrReadOnly = errors.New("read only")
	ErrOp       = errors.New("operation not supported")
	ErrName     = errors.New("bad file or directory name")
	ErrNotFound = errors.New("not found")
)
