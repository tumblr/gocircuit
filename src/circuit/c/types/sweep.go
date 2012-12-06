package types

import (
	"circuit/c/errors"
	"go/ast"
	"go/token"
)

func (tt *TypeTable) AddPackage(fset *token.FileSet, pkgPath string, pkg *ast.Package) error {
	for _, file := range pkg.Files {
		if err := tt.addFile(fset, pkgPath, file); err != nil {
			return err
		}
	}
	return nil
}

func (tt *TypeTable) addFile(fset *token.FileSet, pkgPath string, f *ast.File) error {
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
				if err := tt.addTypeSpec(fset, pkgPath, spec.(*ast.TypeSpec)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (tt *TypeTable) addTypeSpec(fset *token.FileSet, pkgPath string, spec *ast.TypeSpec) error {
	t := &Type{
		Name:    spec.Name.Name,
		PkgPath: pkgPath,
		Spec:    spec,
	}
	if _, ok := tt.types[t.FullName()]; ok {
		return errors.NewSource(fset, spec.Name.NamePos, "type %s already defined", t.FullName())
	}
	tt.types[t.FullName()] = t
	return nil
}

/*
func (tt *TypeTable) linkType(…) … {
	…
	switch q := spec.Type.(type) {
	// Built-in types or references to other types in this package
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
		XX // Slice or array kind?
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
		return nil, errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return r
}
*/
