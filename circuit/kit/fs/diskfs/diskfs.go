// Package localfs exposes a local root directory as a read-only file system
package diskfs

import (
	"errors"
	"os"
	"path"
	fspkg "circuit/kit/fs"
)

var ErrNotDir = errors.New("not a directory")

// FS is a proxy to an isolated subtree of the local file system
type FS struct {
	root     string
	readonly bool
}

func Mount(root string, readonly bool) (*FS, error) {
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, ErrNotDir
	}
	return &FS{
		root:     root,
		readonly: readonly,
	}, nil
}

func (fs *FS) Open(name string) (fspkg.File, error) {
	file, err := os.Open(fs.abs(name))
	if err != nil {
		return nil, err
	}
	return newFile(fs, file), nil
}

func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (fspkg.File, error) {
	// TODO: More rigorous mode and perm checking should happen here
	file, err := os.OpenFile(fs.abs(name), flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(fs, file), nil
}

func (fs *FS) Create(name string) (fspkg.File, error) {
	if fs.readonly {
		return nil, fspkg.ErrReadOnly
	}
	file, err := os.Create(fs.abs(name))
	if err != nil {
		return nil, err
	}
	return newFile(fs, file), nil
}

func (fs *FS) Remove(name string) error {
	if fs.readonly {
		return fspkg.ErrReadOnly
	}
	return os.Remove(fs.abs(name))
}

func (fs *FS) Rename(oldname, newname string) error {
	if fs.readonly {
		return fspkg.ErrReadOnly
	}
	return os.Rename(fs.abs(oldname), fs.abs(newname))
}

func (fs *FS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(fs.abs(name))
}

func (fs *FS) Mkdir(name string) error {
	if fs.readonly {
		return fspkg.ErrReadOnly
	}
	return os.Mkdir(fs.abs(name), 0700)
}

func (fs *FS) MkdirAll(name string) error {
	if fs.readonly {
		return fspkg.ErrReadOnly
	}
	return os.MkdirAll(fs.abs(name), 0700)
}

func (fs *FS) IsReadOnly() bool {
	return fs.readonly
}

func (fs *FS) abs(name string) string {
	return path.Join(fs.root, name)
}

// File represents a proxy to a local file or directory
type File struct {
	fs   *FS
	file *os.File
}

func newFile(fs *FS, file *os.File) *File {
	return &File{fs, file}
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return f.file.Readdir(count)
}

func (f *File) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *File) Truncate(size int64) error {
	if f.fs.IsReadOnly() {
		return fspkg.ErrReadOnly
	}
	return f.file.Truncate(size)
}

func (f *File) Write(q []byte) (int, error) {
	if f.fs.IsReadOnly() {
		return 0, fspkg.ErrReadOnly
	}
	return f.file.Write(q)
}

func (f *File) Sync() error {
	if f.fs.IsReadOnly() {
		return fspkg.ErrReadOnly
	}
	return f.file.Sync()
}
