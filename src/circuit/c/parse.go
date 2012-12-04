package c

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func filterGo(fi os.FileInfo) bool {
	n := fi.Name()
	return len(n) > 0 && strings.HasSuffix(n, ".go") && n[0] != '_'
}

// Skeleton holds a set of parsed packages and their common file set
type Skeleton struct {
	FileSet *token.FileSet
	Pkgs    map[string]*ast.Package
}

// ParsePkg parses a package and returns the result in a new Skeleton
func (b *Build) ParsePkg(pkg string, mode parser.Mode) (ps *Skeleton, err error) {
	ps = &Skeleton{}
	ps.FileSet = token.NewFileSet()
	
	_, pkgpath, err := b.GoPaths.FindPkg(pkg)
	if err != nil {
		return nil, err
	}

	if ps.Pkgs, err = parser.ParseDir(ps.FileSet, pkgpath, filterGo, mode); err != nil {
		return nil, err
	}
	return ps, nil
}
