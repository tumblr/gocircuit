package file

import (
	"os"
	"time"
)

type FileInfo struct {
	SaveName    string
	SaveSize    int64
	SaveMode    os.FileMode
	SaveModTime time.Time
	SaveIsDir   bool
	SaveSys     interface{}
}

func NewFileInfoOS(fi os.FileInfo) *FileInfo {
	return &FileInfo{
		SaveName:    fi.Name(),
		SaveSize:    fi.Size(),
		SaveMode:    fi.Mode(),
		SaveModTime: fi.ModTime(),
		SaveIsDir:   fi.IsDir(),
		SaveSys:     fi.Sys(),
	}
}

func (fi *FileInfo) Name() string {
	return fi.SaveName
}

func (fi *FileInfo) Size() int64 {
	return fi.SaveSize
}

func (fi *FileInfo) Mode() os.FileMode {
	return fi.SaveMode
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.SaveModTime
}

func (fi *FileInfo) IsDir() bool {
	return fi.SaveIsDir
}

func (fi *FileInfo) Sys() interface{} {
	return fi.SaveSys
}
