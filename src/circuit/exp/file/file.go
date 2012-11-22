// Package file provides ways to pass open files to across circuit runtimes
package file

import (
	"os"
	"runtime"
	"circuit/use/circuit"
)

// NewFileClient consumes a circuit pointer, backed by a FileServer on a remote worker, and
// returns a local proxy object with convinient access methods
func NewFileClient(x circuit.X) *FileClient {
	return &FileClient{X: x}
}

type FileClient struct {
	circuit.X
}

func (fcli *FileClient) Close() error {
	return fcli.Call("Close")[0].(error)
}

func (fcli *FileClient) Stat() (os.FileInfo, error) {
	r := fcli.Call("Stat")
	return r[0].(os.FileInfo), r[1].(error)
}

func (fcli *FileClient) Readdir(count int) ([]os.FileInfo, error) {
	r := fcli.Call("Readdir", count)
	return r[0].([]os.FileInfo), r[1].(error)
}

func (fcli *FileClient) Read(p []byte) (int, error) {
	r := fcli.Call("Read", len(p))
	q, err := r[0].([]byte), r[1].(error)
	if len(q) > len(p) {
		panic("corrupt file server")
	}
	copy(p, q)
	return len(q), err
}

func (fcli *FileClient) Seek(offset int64, whence int) (int64, error) {
	r := fcli.Call("Seek", offset, whence)
	return r[0].(int64), r[1].(error)
}

func (fcli *FileClient) Truncate(size int64) error {
	return fcli.Call("Truncate", size)[0].(error)
}

func (fcli *FileClient) Write(p []byte) (int, error) {
	r := fcli.Call("Write", p)
	return r[0].(int), r[1].(error)
}

func (fcli *FileClient) Sync() error {
	return fcli.Call("Sync")[0].(error)
}

// NewFileServer returns a file object which can be passed across runtimes.
// It makes sure to close the file if the no more references to the object remain in the circtui.
func NewFileServer(f *os.File) *FileServer {
	fsrv := &FileServer{f: f}
	runtime.SetFinalizer(fsrv, func(fsrv_ *FileServer) {
		fsrv.f.Close()
	})
	return fsrv
}

type FileServer struct {
	f *os.File
}

func init() {
	circuit.RegisterType(&FileServer{})
}

func (fsrv *FileServer) Close() error {
	return fsrv.f.Close()
}

func (fsrv *FileServer) Stat() (os.FileInfo, error) {
	fi, err := fsrv.f.Stat()
	return NewFileInfoOS(fi), err
}

func (fsrv *FileServer) Readdir(count int) ([]os.FileInfo, error) {
	ff, err := fsrv.f.Readdir(count)
	for i, f := range ff {
		ff[i] = NewFileInfoOS(f)
	}
	return ff, err
}

func (fsrv *FileServer) Read(n int) ([]byte, error) {
	p := make([]byte, min(n, 1e4))
	m, err := fsrv.f.Read(p)
	return p[:m], err
}

func min (x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (fsrv *FileServer) Seek(offset int64, whence int) (int64, error) {
	return fsrv.f.Seek(offset, whence)
}

func (fsrv *FileServer) Truncate(size int64) error {
	return fsrv.f.Truncate(size)
}

func (fsrv *FileServer) Write(p []byte) (int, error) {
	return fsrv.f.Write(p)
}

func (fsrv *FileServer) Sync() error {
	return fsrv.f.Sync()
}
