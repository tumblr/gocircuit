package c

import (
	"go/ast"
	"go/printer"
	"path"
)

// flush writes out all compiled and transformed packages to their location
// inside the compilation jail
func (b *Build) flush() error {
	for pkgPath, pkg := range b.pkgs {
		for fileName, fileAST := range pkg.Files {
			_, fileName = path.Split(fileName)
			if err := b.flushFile(pkgPath, fileName, fileAST); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Build) flushFile(pkgPath, fileName string, file *ast.File) error {
	f, err := b.jail.CreateSrcFile(pkgPath, fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, b.fileSet, file)
}
