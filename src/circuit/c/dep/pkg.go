package pkg

import (
	"path"
)

type Pkg struct {
	Path string    // Package path within root GOPATH, e.g. circuit/use/anchorfs
	Name string    // Name of package,                 e.g. anchorfs
}

func NewPkgImport(importPath string) *Pkg {
	_, name := path.Split(importPath)
	return &Pkg{
		Path: importPath,
		Name: name,
	}
}
