package c

import (
	"circuit/c/types"
	"go/ast"
)

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

	for _, pkgSrc := range b.pkgs {

		// Create source file for registrations in this package
		pkgName := pkgSrc.Name()
		/*tokenFile, astFile :=*/ pkgSrc.AddFile(pkgName, pkgName + "_circuit.go") 

		for _/*pkgName*/, pkg := range pkgSrc.Pkgs {
			
			// Compile interface types in package
			if err := types.CompilePkg(pkgSrc.FileSet, pkg, func (spec *ast.TypeSpec) error {
					return nil
				}); err != nil {
				return err
			}

			// Write interface registrations
			// …
		}
	}
	return nil
}

//func registerInterface(spec *ast.TypeSpec) error {
//}
