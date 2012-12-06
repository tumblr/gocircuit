package c

import (
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"strings"
)

func filterGoNoTest(fi os.FileInfo) bool {
	n := fi.Name()
	return len(n) > 0 && strings.HasSuffix(n, ".go") && n[0] != '_' && strings.Index(n, "_test.go") < 0
}

// ParsePkg parses package pkg, using FileSet fset
func ParsePkg(l *Layout, fset *token.FileSet, pkgPath string, mode parser.Mode) (pkgs map[string]*ast.Package, err error) {
	if _, pkgPath, err = l.FindPkg(pkgPath, true); err != nil {
		return nil, err
	}

	if pkgs, err = parser.ParseDir(fset, pkgPath, filterGoNoTest, mode); err != nil {
		return nil, err
	}
	return pkgs, nil
}
