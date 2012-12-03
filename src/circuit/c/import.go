package c

import (
	"go/ast"
	"path"
	"strconv"
)

func packageImports(pkg *ast.Package) map[string]struct{} {
	imprts := make(map[string]struct{}) 
	for _, file := range pkg.Files {
		for _, impSpec := range file.Imports {
			imprts[newFileImport(impSpec).Path] = struct{}{}
		}
	}
	return imprts
}

type fileImport struct {
	Name string  // local import name
	Path string  // package full path
}

func newFileImport(spec *ast.ImportSpec) *fileImport {
	var err error
	fimp := &fileImport{}
	if fimp.Path, err = strconv.Unquote(spec.Path.Value); err != nil {
		panic(err)
	}
	if spec.Name == nil {
		_, fimp.Name = path.Split(fimp.Path)
	} else {
		fimp.Name = spec.Name.Name
	}
	return fimp
}

// fileImports returns a slice of all packages imported by file
func fileImports(file *ast.File) []*fileImport {
	imprts := make([]*fileImport, len(file.Imports))
	for i, impSpec := range file.Imports {
		imprts[i] = newFileImport(impSpec)
	}
	return imprts
}
