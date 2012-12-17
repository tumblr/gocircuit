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

// PkgPath returns the absolute path of package pkg within the jail
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
