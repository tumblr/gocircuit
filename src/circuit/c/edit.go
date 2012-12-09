package c

import (
	"go/ast"
	"go/token"
)

// name is the just file name (excluding any path prefix)
func AddFile(fset *token.FileSet, pkg *ast.Package, name string) (*token.File, *ast.File) {

	// Add file to file set
	ff := fset.AddFile(name, fset.Base(), 1)
	if pkg.Files == nil {
		pkg.Files = make(map[string]*ast.File)
	}

	// Does source file already exist in package?
	if _, present := pkg.Files[name]; present {
		panic("file already added")
	}

	// Add file to AST
	file := &ast.File{
		Package: ff.Pos(0),
		Name:    &ast.Ident{
			NamePos: ff.Pos(0),
			Name:    pkg.Name,
		},
	}
	pkg.Files[name] = file

	return ff, file
}
