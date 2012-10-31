package anchorfs

import (
	"path"
	"strings"
	"time"
	"tumblr/circuit/use/lang"
)

var (
	ErrName     = lang.NewError("anchor name")
	ErrNotFound = lang.NewError("not found")
)

// fs represents an anchor file system
type fs interface {
	CreateFile(string, lang.Addr) error
	OpenFile(string) (File, error)
	OpenDir(string) (Dir, error)
	Created() []string
}

// Dir is the interface for a directory of workers in the anchor file system
type Dir interface {
	Name() string
	Dirs() ([]string, error)

	Files()                                            (rev int64, workers map[lang.RuntimeID]File, err error)
	Change(sinceRev int64)                             (rev int64, workers map[lang.RuntimeID]File, err error)
	ChangeExpire(sinceRev int64, expire time.Duration) (rev int64, workers map[lang.RuntimeID]File, err error)

	OpenFile(lang.RuntimeID) (File, error)
	OpenDir(string) (Dir, error)
}

// File ...
type File interface {
	Owner() lang.Addr
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
		if _, err := lang.ParseRuntimeID(part); err == nil {
			return nil, "", ErrName
		}
	}
	return parts, "/" + path.Join(parts...), nil
}
