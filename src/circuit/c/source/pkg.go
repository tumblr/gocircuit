package source

import (
	"go/ast"
	"go/token"
	"path"
	"strings"
)

// Pkg captures a parsed Go source package
type Pkg struct {
	FileSet *token.FileSet          // File names are relative to SrcDir
	SrcDir  string                  // SrcDir/PkgPath = absolute local path to package directory
	PkgPath string                  // Package import path
	PkgAST  map[string]*ast.Package // Package name to package AST
	FileAST map[string]*ast.File
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

func (p *Pkg) GetPkg(name string) *ast.Package {
	return p.PkgAST[name]
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

func (p *Pkg) AddPkg(name string) *ast.Package {
	pkg, ok := p.PkgAST[name]
	if !ok {
		pkg = &ast.Package{
			Name:  name,
			Files: make(map[string]*ast.File),
		}
		p.PkgAST[name] = pkg
	}
	return pkg
}

func (p *Pkg) AddFile(pkgName, fileName string) *ast.File {
	if strings.Index(fileName, "/") >= 0 {
		panic("not a filename")
	}
	pkg := p.AddPkg(pkgName)

	filePath := path.Join(p.PkgPath, fileName)
	f, ok := pkg.Files[filePath]
	if !ok {
		ff := p.FileSet.AddFile(filePath, p.FileSet.Base(), 1)
		pos := ff.Pos(0)
		f = &ast.File{
			Package: pos,
			Name:    &ast.Ident{Name: pkgName},
		}
		pkg.Files[filePath] = f
	}

	return f
}
