package c

import (
	"circuit/c/types"
	"go/ast"
)

// TODO: copy type declarations if not from main package

// For every package path, and every parsed package (name) inside of it, fish
// out all public types and register them with the circuit type system in a
// single new source file whose package is named after the package path.
//
// The rationale is that we want circuit workers to have types from packages as
// well as executables linked in. This addresses the situation when an entire
// circuit app is implemented in a "main" package and the corresponding circuit
// worker must be aware of the types declared in that "main" package.
//
func (b *Build) transformRegisterInterfaces() error {

	// For every package directory
	for _, pkgSrc := range b.pkgs {

		// Create source file for registrations in this package
		pkgName := pkgSrc.Name()
		_, astFile := pkgSrc.AddFile(pkgName, pkgName + "_circuit.go") 
		AddImport(astFile, "circuit/use/circuit")

		// For every package name defined in the package directory
		for _, pkg := range pkgSrc.Pkgs {

			var typeNames []string
			
			// Compile interface types in package
			if err := types.CompilePkg(pkgSrc.FileSet, pkg, func (spec *ast.TypeSpec) error {
					typeNames = append(typeNames, spec.Name.Name)
					return nil
				}); err != nil {
				return err
			}

			// Write interface registrations
			registerValues(astFile, typeNames)
		}
	}
	return nil
}

func registerValues(file *ast.File, typeNames []string) {

	// Create init function declaration
	fdecl := &ast.FuncDecl{
		Doc:  nil,
		Recv: nil,
		Name: &ast.Ident{Name: "init"},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}
	file.Decls = append(file.Decls, fdecl)

	// Add type registrations to fdecl.Body.List
	for _, name := range typeNames {
		stmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "circuit" },       // Refers to import circuit/use/circuit
					Sel: &ast.Ident{Name: "RegisterValue"},
				},
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.Ident{Name: name},
					},
				},
			},
		}
		fdecl.Body.List = append(fdecl.Body.List, stmt)
	}

}
