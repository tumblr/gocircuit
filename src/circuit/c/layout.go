package c

import (
	"os"
	"path"
	"strings"
)

// Layout describes a Go compilation environment
type Layout struct {
	goRoot        string    // GOROOT directory
	workingGoPath string    // GOPATH of the user Go repo
	goPaths       GoPaths   // All GOPATH paths
}

// NewWorkingLayout creates a new build environment, where the working
// gopath is derived from the current working directory.
func NewWorkingLayout() (*Layout, error) {
	gopath, err := FindWorkingGoPath()
	if err != nil {
		return nil, err
	}
	return &Layout{
		goRoot:        os.Getenv("GOROOT"),
		workingGoPath: gopath,
		goPaths:       GetGoPaths(),
	}, nil
}

// FindPkg returns the first gopath that contains package pkg.
// If includeGoRoot is set, goroot is checked first.
func (l *Layout) FindPkg(pkg string, includeGoRoot bool) (gopath, pkgpath string, err error) {
	if includeGoRoot {
		pkgpath, err = existPkg(l.goRoot, path.Join("pkg", pkg))
		if err == nil {
			return l.goRoot, pkgpath, nil
		}
		if err != ErrNotFound {
			return "", "", err
		}
	}
	return l.goPaths.FindPkg(pkg)
}

// FindWorkingPath returns the first gopath that parents the absolute directory dir.
// If includeGoRoot is set, goroot is checked first.
func (l *Layout) FindWorkingPath(dir string, includeGoRoot bool) (gopath string, err error) {
	if includeGoRoot {
		if strings.HasPrefix(dir, path.Join(l.goRoot, "src")) {
			return l.goRoot, nil
		}
	}
	return l.goPaths.FindWorkingPath(dir)
}
