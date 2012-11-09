package fs

import (
	"os"
	"time"
)

type FileInfo struct {
	XName    string
	XSize    int64
	XMode    os.FileMode
	XModTime time.Time
	XIsDir   bool
}

func (fi *FileInfo) Name() string {
	return fi.XName
}

func (fi *FileInfo) Size() int64 {
	return fi.XSize
}

func (fi *FileInfo) Mode() os.FileMode {
	return fi.XMode
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.XModTime
}

func (fi *FileInfo) IsDir() bool {
	return fi.XIsDir
}

func (fi *FileInfo) Sys() interface{} {
	return nil
}
