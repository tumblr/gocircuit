package fs

import "net/http"

func HTTPFileSystem(fs FS) http.FileSystem {
	return httpFS{fs}
}

type httpFS struct {
	fs FS
}

func (httpfs httpFS) Open(name string) (http.File, error) {
	return httpfs.fs.Open(name)
}
