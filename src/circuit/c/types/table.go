package types

import (
	"circuit/c/errors"
	"sort"
)

type TypeTable struct {
	types map[string]*Type	           // Fully-qualified type name to type structure
	pkgs  map[string]map[string]*Type  // Package path to type name to type structure
}

func NewTypeTable() *TypeTable {
	return &TypeTable{
		types: make(map[string]*Type),
		pkgs:  make(map[string]map[string]*Type),
	}
}

func (tt *TypeTable) ListFullNames() []string {
	var pp []string
	for pkgPath, _ := range tt.types {
		pp = append(pp, pkgPath)
	}
	sort.Strings(pp)
	return pp
}

func (tt *TypeTable) addType(t *Type) error {

	// Add type to global type map
	if _, ok := tt.types[t.FullName()]; ok {
		return errors.NewSource(t.FileSet, t.Spec.Name.NamePos, "type %s already defined", t.FullName())
	}
	tt.types[t.FullName()] = t

	// Add type to per-package structure
	pkgMap, ok := tt.pkgs[t.PkgPath] 
	if !ok {
		pkgMap = make(map[string]*Type)
		tt.pkgs[t.PkgPath] = pkgMap
	}
	pkgMap[t.Name] = t

	return nil
}

// PkgTypes returns a map from type name to type structure of all types declared in pkgPath
func (tt *TypeTable) PkgTypes(pkgPath string) map[string]*Type {
	pkgTypes, ok := tt.pkgs[pkgPath]
	if !ok {
		return nil
	}
	return pkgTypes
}
