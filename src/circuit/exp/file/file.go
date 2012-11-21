// Package file provides ways to pass open files to across circuit runtimes
package clisrv

import (
	"os"
	"runtime"
	"circuit/use/circuit"
)

// NewFileClient consumes a circuit pointer, backed by a FileServer on a remote worker, and
// returns a local proxy object with convinient access methods
func NewFileClient(x circuit.X) FileClient {
	panic("not implemented")	
}

type FileClient circuit.X

// NewFile returns a file object which can be passed across runtimes.
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
	return fsrv.f.Stat()
}

func (fsrv *FileServer) Readdir(count int) ([]os.FileInfo, error) {
	return fsrv.f.Readdir(count)
}

func (fsrv *FileServer) Read(p []byte) (int, error) {
	panic("won't work until circuit supports passing mutable slices")
	return fsrv.f.Read(p)
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
