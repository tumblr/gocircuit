package types

import (
	"sort"
)


type TypeTable struct {
	types map[string]*Type
}

func NewTypeTable() *TypeTable {
	return &TypeTable{
		types: make(map[string]*Type),
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
