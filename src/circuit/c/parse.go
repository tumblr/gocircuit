package c

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func filterGoNoTest(fi os.FileInfo) bool {
	n := fi.Name()
	return len(n) > 0 && strings.HasSuffix(n, ".go") && n[0] != '_' && strings.Index(n, "_test.go") < 0
}

// Skeleton holds a set of parsed packages and their common file set
type Skeleton struct {
	FileSet *token.FileSet
	Pkgs    map[string]*ast.Package
}

// ParsePkg parses a package and returns the result in a new Skeleton
func (l *Layout) ParsePkg(pkg string, mode parser.Mode) (ps *Skeleton, err error) {
	ps = &Skeleton{}
	ps.FileSet = token.NewFileSet()
	
	_, pkgpath, err := l.FindPkg(pkg, true)
	if err != nil {
		return nil, err
	}

	if ps.Pkgs, err = parser.ParseDir(ps.FileSet, pkgpath, filterGoNoTest, mode); err != nil {
		return nil, err
	}
	return ps, nil
}
