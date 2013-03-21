// Package file provides ways to pass open files to across circuit runtimes
package file

import (
	"circuit/use/circuit"
	"encoding/gob"
	"os"
	"runtime"
)

func init() {
	gob.Register(&os.PathError{})
}

// NewFileClient consumes a circuit pointer, backed by a FileServer on a remote worker, and
// returns a local proxy object with convinient access methods
func NewFileClient(x circuit.X) *FileClient {
	return &FileClient{X: x}
}

type FileClient struct {
	circuit.X
}

func asError(x interface{}) error {
	if x == nil {
		return nil
	}
	return x.(error)
}

func asFileInfo(x interface{}) os.FileInfo {
	if x == nil {
		return nil
	}
	return x.(os.FileInfo)
}

func asFileInfoSlice(x interface{}) []os.FileInfo {
	if x == nil {
		return nil
	}
	return x.([]os.FileInfo)
}

func asBytes(x interface{}) []byte {
	if x == nil {
		return nil
	}
	return x.([]byte)
}

func fileRecover(pe *error) {
	if p := recover(); p != nil {
		*pe = circuit.NewError("server died")
	}
}

func (fcli *FileClient) Close() (err error) {
	defer fileRecover(&err)

	return asError(fcli.Call("Close")[0])
}

func (fcli *FileClient) Stat() (_ os.FileInfo, err error) {
	defer fileRecover(&err)

	r := fcli.Call("Stat")
	return asFileInfo(r[0]), asError(r[1])
}

func (fcli *FileClient) Readdir(count int) (_ []os.FileInfo, err error) {
	defer fileRecover(&err)

	r := fcli.Call("Readdir", count)
	return asFileInfoSlice(r[0]), asError(r[1])
}

func (fcli *FileClient) Read(p []byte) (_ int, err error) {
	defer fileRecover(&err)

	r := fcli.Call("Read", len(p))
	q, err := asBytes(r[0]), asError(r[1])
	if len(q) > len(p) {
		panic("corrupt file server")
	}
	copy(p, q)
	return len(q), err
}

func (fcli *FileClient) Seek(offset int64, whence int) (_ int64, err error) {
	defer fileRecover(&err)

	r := fcli.Call("Seek", offset, whence)
	return r[0].(int64), asError(r[1])
}

func (fcli *FileClient) Truncate(size int64) (err error) {
	defer fileRecover(&err)

	return asError(fcli.Call("Truncate", size)[0])
}

func (fcli *FileClient) Write(p []byte) (_ int, err error) {
	defer fileRecover(&err)

	r := fcli.Call("Write", p)
	return r[0].(int), asError(r[1])
}

func (fcli *FileClient) Sync() (err error) {
	defer fileRecover(&err)

	return asError(fcli.Call("Sync")[0])
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
	circuit.RegisterValue(&FileServer{})
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

func min(x, y int) int {
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
