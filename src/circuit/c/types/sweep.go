package types

import (
	"go/ast"
	"go/token"
	"reflect"
)

type TypeTable struct {
	types map[string]*Type
}

type Type struct {
	Spec  *ast.TypeSpec
	Name  string
	Kind  reflect.Kind
	Elem  *Type
	Procs []*Proc
}

type Proc struct {
	AST *ast.Node
}

func NewTypeTable() *TypeTable {
	return &TypeTable{
		types: make(map[string]*Type),
	}
}

func (tt *TypeTable) ParseFile(fset *ast.FileSet, f *ast.File) error {
	for _, decl := range f.Decls {
		switch q := decl.(type) {
		// GenDecl captures a single or multi-type declaration block, e.g.:
		//	type T0 … 
		//	type (
		//		T1 …
		//		T2 …
		//	)
		case *ast.GenDecl:
			if q.Tok != token.TYPE {
				break
			}
			for _, spec := range q.Specs {
				typ, err := ParseTypeSpec(fset, spec.(*ast.TypeSpec))
				if err != nil {
					return err
				}
			}
		}
	}
}

func ParseTypeSpec(fset *ast.FileSet, spec *ast.TypeSpec) (*Type, error) {
	r := &Type{
		Name: spec.Name.Ident.Name,
		Spec: spec,
	}
	
	switch q := spec.Type.(type) {
	case *ast.Ident:
		?
	case *ast.ParenExpr:
		?
	case *ast.SelectorExpr:
		?
	case *ast.StarExpr:
		// r.Elem will be filled in during a follow up sweep of all types
		r.Kind = reflect.Ptr
	case *ast.ArrayType:
		r.Kind = reflect.Array
	case *ast.ChanType:
		r.Kind = reflect.Chan
	case *ast.FuncType:
		r.Kind = reflect.Func
	case *ast.InterfaceType:
		r.Kind = reflect.Interface
	case *ast.MapType:
		r.Kind = reflect.Map
	case *ast.StructType:
		r.Kind = reflect.Struct
	default:
		return nil, NewError(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return r
}
