package c

import (
	"go/ast"
	"go/token"
)

// GetPkg …
func (p *PkgSrc) GetPkg(name string) *ast.Package {
	pkg, ok := p.Pkgs[name]
	if !ok {
		pkg = &ast.Package{
			Name:  name,
			Files: make(map[string]*ast.File),
		}
		p.Pkgs[name] = pkg
	}
	return pkg
}

// name is the just file name (excluding any path prefix)
func (p *PkgSrc) AddFile(pkgName, fileName string) (*token.File, *ast.File) {


	// Add file to file set
	ff := p.FileSet.AddFile(fileName, p.FileSet.Base(), 1)

	// Fetch package AST
	pkg := p.GetPkg(pkgName)
	if pkg.Files == nil {
		pkg.Files = make(map[string]*ast.File)
	}

	// Does source file already exist in package?
	if _, present := pkg.Files[fileName]; present {
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
	pkg.Files[fileName] = file

	return ff, file
}
