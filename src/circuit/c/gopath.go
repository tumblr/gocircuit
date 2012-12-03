package c

import (
	"errors"
	"os"
	"path"
	"sort"
	"strings"
)


// GetWorkingGoPath returns the most specific GOPATH for the current working directory
func GetWorkingGoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return GetGoPath(wd)
}

// GetGoPath returns the most specific GOPATH for the given directory
func GetGoPath(dir string) (string, error) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	for i, gp := range gopaths {
		gopaths[i] = path.Clean(gp)
	}
	sort.Sort(descendingLenStrings(gopaths))
	dir = path.Clean(dir)
	for _, gp := range gopaths {
		if strings.HasPrefix(dir, gp) {
			return gp, nil
		}
	}
	return "", errors.New("gopath not found")
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
