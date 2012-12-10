package types

import (
	"go/ast"
	"path"
	"reflect"
)

type Type struct {

	// Sweep 1
	Spec    *ast.TypeSpec
	Name    string
	PkgPath string

	// Sweep 2
	Kind    reflect.Kind
	Elem    *Type
	Procs   []*Proc
}

func (t *Type) FullName() string {
	return path.Join(t.PkgPath, t.Name)
}
