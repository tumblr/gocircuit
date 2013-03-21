// Package anchorfs exposes the programming interface for accessing the anchor file system
package anchorfs

import (
	"circuit/use/circuit"
	"path"
	"strings"
	"time"
)

var (
	ErrName     = circuit.NewError("anchor name")
	ErrNotFound = circuit.NewError("not found")
)

// fs represents an anchor file system
type fs interface {
	CreateFile(string, circuit.Addr) error
	OpenFile(string) (File, error)
	OpenDir(string) (Dir, error)
	Created() []string
}

// Dir is the interface for a directory of workers in the anchor file system
type Dir interface {
	Name() string
	Dirs() ([]string, error)

	Files() (rev int64, workers map[circuit.WorkerID]File, err error)
	Change(sinceRev int64) (rev int64, workers map[circuit.WorkerID]File, err error)
	ChangeExpire(sinceRev int64, expire time.Duration) (rev int64, workers map[circuit.WorkerID]File, err error)

	OpenFile(circuit.WorkerID) (File, error)
	OpenDir(string) (Dir, error)
}

// File ...
type File interface {
	Owner() circuit.Addr
}

// Sanitizer ensures that anchor is a valid anchor path in the fs
// and returns its parts
func Sanitize(anchor string) ([]string, string, error) {
	anchor = path.Clean(anchor)
	if len(anchor) == 0 || anchor[0] != '/' {
		return nil, "", ErrName
	}
	parts := strings.Split(anchor[1:], "/")
	for _, part := range parts {
		if _, err := circuit.ParseWorkerID(part); err == nil {
			return nil, "", ErrName
		}
	}
	return parts, "/" + path.Join(parts...), nil
}
