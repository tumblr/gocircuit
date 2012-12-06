package c

import (
	"circuit/c/errors"
	"os"
	"path"
	"sort"
	"strings"
)

// GoPaths is a structure providing query facilities for a GOPATH environment
type GoPaths []string

// NewGoPaths creates a GoPaths structure from the colon-separated list of paths gopathlist
func NewGoPaths(gopathlist string) GoPaths {
	gopaths := strings.Split(gopathlist, ":")
	for i, gp := range gopaths {
		gopaths[i] = path.Clean(gp)
	}
	return GoPaths(gopaths)
}

// FindPkg looks for pkg in each gopath in order of appearance.
// If found, it returns the gopath as well as the absolute package path.
func (gopaths GoPaths) FindPkg(pkg string) (gopath, pkgpath string, err error) {
	for _, gp := range gopaths {
		pkgpath, err = existPkg(gp, pkg)
		if err == errors.ErrNotFound {
			continue
		}
		if err != nil {
			return "", "", err
		}
		return gp, pkgpath, nil
	}
	return "", "", errors.ErrNotFound
}

func existPkg(gopath, pkg string) (pkgpath string, err error) {
	pkgpath = path.Join(gopath, "src", pkg)
	fi, err := os.Stat(pkgpath)
	if os.IsNotExist(err) {
		return "", errors.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", errors.ErrNotFound
	}
	return pkgpath, nil
}

func (gopaths GoPaths) FindWorkingPath(dir string) (string, error) {
	order := make([]string, len(gopaths))
	copy(order, gopaths)
	sort.Sort(descendingLenStrings(order))

	dir = path.Clean(dir)
	for _, gp := range order {
		if strings.HasPrefix(dir, path.Join(gp, "src")) {
			return gp, nil
		}
	}

	return "", errors.New("working gopath not found")
}

// GetGoPaths returns a gopath structure for the current environment GOPATH
func GetGoPaths() GoPaths {
	return NewGoPaths(os.Getenv("GOPATH"))
}

// FindGoPath returns the most specific GOPATH for the given directory
func FindGoPath(dir string) (string, error) {
	return GetGoPaths().FindWorkingPath(dir)
}

// FindWorkingGoPath returns the most specific GOPATH for the current working directory
func FindWorkingGoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return FindGoPath(wd)
}

type descendingLenStrings []string

func (t descendingLenStrings) Len() int {
	return len(t)
}

func (t descendingLenStrings) Less(i, j int) bool {
	if len(t[i]) == len(t[j]) {
		return t[i] < t[j]
	}
	return len(t[i]) > len(t[j])
}

func (t descendingLenStrings) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
