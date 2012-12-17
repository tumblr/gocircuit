package source

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path"
)

// Flush writes out all compiled and transformed packages to their location
// inside the compilation jail
func (s *Source) Flush() error {
	// For every Pkg
	for pkgPath, _ := range s.pkg {
		if err := s.FlushPkg(pkgPath); err != nil {
			return err
		}
	}
	return nil
}

func (s *Source) FlushPkg(pkgPath string) error {

	pkg := s.GetPkg(pkgPath)

	// For every ast.Package
	for _, pkgAST := range pkg.PkgAST {
		// For every ast.File
		for filePath, fileAST := range pkgAST.Files {
			_, fileName := path.Split(filePath)
			f, err := s.jail.CreateSourceFile(pkgPath, fileName)
			if err != nil {
				return err
			}
			if err := flushFile(f, pkg.FileSet, fileAST); err != nil {
				return err
			}
		}
	}
	return nil
}

func flushFile(f *os.File, fileSet *token.FileSet, file *ast.File) error {
	defer f.Close()
	return printer.Fprint(f, fileSet, file)
}
