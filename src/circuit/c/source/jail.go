package source

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

// AbsPkgPath returns the absolute local path of package pkg within the jail
func (j *Jail) AbsPkgPath(pkgPath string) string {
	return path.Join(j.root, "src", pkgPath)
}

func (j *Jail) MakePkgDir(pkgPaths ...string) error {
	for _, pkgPath := range pkgPaths {
		if err := os.MkdirAll(j.AbsPkgPath(pkgPath), 0700); err != nil {
			return err
		}
	}
	return nil
}

func (j *Jail) CreateSourceFile(pkgPath, fileName string) (*os.File, error) {
	absPath := j.AbsPkgPath(pkgPath)
	if err := os.MkdirAll(absPath, 0770); err != nil {
		return nil, err
	}
	f, err := os.Create(path.Join(absPath, fileName))
	if err != nil {
		return nil, err
	}
	return f, nil
}
