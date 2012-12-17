package source

import (
	"go/ast"
	"go/token"
	"path"
)

// Pkg captures a parsed Go source package
type Pkg struct {
	FileSet  *token.FileSet           // File names are relative to SrcDir
	SrcDir   string                   // SrcDir/PkgPath = absolute local path to package directory
	PkgPath  string                   // Package import path
	PkgAST   map[string]*ast.Package  // Package name to package AST
	FileAST  map[string]*ast.File
}

func (p *Pkg) link() {
	p.FileAST = make(map[string]*ast.File)
	for _, pkgAST := range p.PkgAST {
		for n, f := range pkgAST.Files {
			if _, ok := p.FileAST[n]; ok {
				panic("file in two packages")
			}
			p.FileAST[n] = f
		}
	}
}

func (p *Pkg) LibPkg() *ast.Package {
	name := p.Name()
	for pkgName, pkg := range p.PkgAST {
		if pkgName == name {
			return pkg
		}
	}
	return nil
}

func (p *Pkg) MainPkg() *ast.Package {
	for pkgName, pkg := range p.PkgAST {
		if pkgName == "main" {
			return pkg
		}
	}
	return nil
}

func (p *Pkg) Name() string {
	_, name := path.Split(p.PkgPath)
	return name
}

func (p *Pkg) AddFile(pkgName, fileName string) *ast.File {

	// Make package ast if not there
	pkg, ok := p.PkgAST[pkgName]
	if !ok {
		pkg = &ast.Package{
			Name:    pkgName,
			Scope:   nil,
			Imports: nil,
			Files:   make(map[string]*ast.File),
		}
		p.PkgAST[pkgName] = pkg
	}

	// If file already exists, return it
	f, ok := pkg.Files[fileName]
	if !ok {
		ff := p.FileSet.AddFile(path.Join(p.PkgPath, fileName), p.FileSet.Base(), 1)
		pos := ff.Pos(0)
		f = &ast.File{
			Package:   pos,
			Name:      &ast.Ident{Name: pkgName},
		}
		pkg.Files[fileName] = f
	}

	return f
}
