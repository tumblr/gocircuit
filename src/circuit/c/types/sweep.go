package types

import (
	"go/ast"
	"go/token"
)

func CompilePkg(fset *token.FileSet, pkgPath string, pkg *ast.Package, globalNames *GlobalNames) error {
	return VisitPkgTypeSpecs(fset, pkg, func(spec *ast.TypeSpec) error {
		t, err := CompileTypeSpec(pkgPath, spec, fimp??)
		if err != nil {
			return err
		}
		switch q := t.(type) {
		case *Named:
			return globalNames.Add(q)
		}
		return nil
	})
}

type TypeSpecFunc func(typeSpec *ast.TypeSpec) error

// VisitPkgTypeSpecs calls typeSpecFunc for each TypeSpec in package pkg.
func VisitPkgTypeSpecs(fset *token.FileSet, pkg *ast.Package, typeSpecFunc TypeSpecFunc) error {
	for _, file := range pkg.Files {
		if err := VisitFileTypeSpecs(fset, file, typeSpecFunc); err != nil {
			return err
		}
	}
	return nil
}

// VisitFileTypeSpecs calls typeSpecFunc for each TypeSpec in file f.
func VisitFileTypeSpecs(fset *token.FileSet, f *ast.File, typeSpecFunc TypeSpecFunc) error {
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
