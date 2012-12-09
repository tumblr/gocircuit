package c

import (
	"go/ast"
	"go/printer"
	"go/token"
	"path"
)

// flush writes out all compiled and transformed packages to their location
// inside the compilation jail
func (b *Build) flush() error {
	for pkgPath, pkgSrc := range b.pkgs {
		for _, pkg := range pkgSrc.Pkgs {
			for fileName, fileAST := range pkg.Files {
				_, fileName = path.Split(fileName)
				if err := b.flushFile(pkgPath, fileName, pkgSrc.FileSet, fileAST); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *Build) flushFile(pkgPath, fileName string, fileSet *token.FileSet, file *ast.File) error {
	f, err := b.jail.CreateSrcFile(pkgPath, fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, fileSet, file)
}
