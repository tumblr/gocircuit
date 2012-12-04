package c

import (
	"os"
	"path"
)

type Jail struct {
	root string
}

func NewJail(root string) (*Jail, error) {
	j := &Jail{root}
	if err := j.mkdirs(); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *Jail) mkdirs() error {
	if err := os.MkdirAll(path.Join(j.root, "src"), 0700); err != nil {
		return err
	}
	return nil
}

// PkgPath returns the absolute path of package pkg within the jail
func (j *Jail) PkgPath(pkg string) string {
	return path.Join(j.root, "src", pkg)
}

func (j *Jail) MakePkgDir(pkgs ...string) error {
	for _, pkg := range pkgs {
		if err := os.MkdirAll(j.PkgPath(pkg), 0700); err != nil {
			return err
		}
	}
	return nil
}

// CreateSrcFile creates an empty writable file name within the package directory of pkg in the jail
func (j *Jail) CreateSrcFile(pkg, name string) (*os.File, error) {
	abs := path.Join(j.PkgPath(pkg), name)
	if err := j.MakePkgDir(pkg); err != nil {
		return nil, err
	}
	return os.Create(abs)
}
