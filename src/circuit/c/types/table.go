package types

import (
	"circuit/c/errors"
	"sort"
)

type GlobalNames struct {
	names map[string]*Named       // Fully-qualified type name to type structure
	pkgs  map[string]PackageNames // Package path to type name to type structure
}

type PackageNames map[string]*Named

func MakeNames() *GlobalNames {
	return &GlobalNames{
		names: make(map[string]*Named),
		pkgs:  make(map[string]PackageNames),
	}
}

func (tt *GlobalNames) ListFullNames() []string {
	var pp []string
	for name, _ := range tt.names {
		pp = append(pp, name)
	}
	sort.Strings(pp)
	return pp
}

// add adds t to the structures for global and per-package lookups
func (tt *GlobalNames) Add(t *Named) error {

	// Add type to global type map
	if _, ok := tt.names[t.FullName()]; ok {
		return errors.New("type %s already defined", t.FullName())
	}
	tt.names[t.FullName()] = t

	// Add type to per-package structure
	pkgMap, ok := tt.pkgs[t.PkgPath]
	if !ok {
		pkgMap = make(map[string]*Named)
		tt.pkgs[t.PkgPath] = pkgMap
	}
	pkgMap[t.Name] = t

	return nil
}

// Pkg returns a map from type name to type structure of all names declared in pkgPath
func (tt *GlobalNames) Pkg(pkgPath string) PackageNames {
	return tt.pkgs[pkgPath]
}
