package util

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// DeriveFileDeps imports returns a slice of all packages that file depends on
func DeriveFileDeps(file *ast.File) []*Pkg {
	?
}

func PkgDeps(pkg string) (map[string]struct{}, error) {
	?
}
