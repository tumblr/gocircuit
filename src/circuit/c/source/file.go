package parse

import (
	"go/ast"
)

type File struct {
	PkgPath string
	AST     *ast.File
}

// LookupImport returns the full package path of the package imported under alias
func (f *File) LookupImport(alias string) string {
	?
}
