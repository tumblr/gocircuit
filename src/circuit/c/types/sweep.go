package types

import (
	"circuit/c/errors"
	"go/ast"
	"go/token"
)

func (tt *TypeTable) AddPkg(fset *token.FileSet, pkgPath string, pkg *ast.Package) error {
	return CompilePkg(fset, pkg, func(spec *ast.TypeSpec) error {
		t := &Type{
			FileSet: fset,
			Name:    spec.Name.Name,
			PkgPath: pkgPath,
			Spec:    spec,
		}
		if _, ok := tt.types[t.FullName()]; ok {
			return errors.NewSource(fset, spec.Name.NamePos, "type %s already defined", t.FullName())
		}
		tt.types[t.FullName()] = t
		return nil
	})
}

type TypeSpecFunc func(typeSpec *ast.TypeSpec) error

func CompilePkg(fset *token.FileSet, pkg *ast.Package, typeSpecFunc TypeSpecFunc) error {
	for _, file := range pkg.Files {
		if err := CompileFile(fset, file, typeSpecFunc); err != nil {
			return err
		}
	}
	return nil
}

func CompileFile(fset *token.FileSet, f *ast.File, typeSpecFunc TypeSpecFunc) error {
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
				if err := typeSpecFunc(spec.(*ast.TypeSpec)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
