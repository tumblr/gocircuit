package c

import (
	"os"
	"path"
	"strings"
)

// Layout describes a Go compilation environment
type Layout struct {
	goRoot        string    // GOROOT directory
	goPaths       GoPaths   // All GOPATH paths
	workingGoPath string    // A distinct GOPATH
}

func NewLayout(goroot string, gopaths GoPaths, working string) *Layout {
	return &Layout{
		goRoot:        goroot,
		goPaths:       gopaths,
		workingGoPath: working,
	}
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

// FindPkg returns the ...
// If includeGoRoot is set, goroot is checked first.
func (l *Layout) FindPkg(pkgPath string, includeGoRoot bool) (srcDir string, err error) {
	if includeGoRoot {
		if err = ExistPkg(path.Join(l.goRoot, "src", "pkg", pkgPath)); err != nil {
			return "", err
		}
		return path.Join(l.goRoot, "src", "pkg"), nil
	}
	return l.goPaths.FindPkg(pkgPath)
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
