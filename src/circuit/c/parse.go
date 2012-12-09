package c

import (
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"path"
	"strings"
)

// PkgSource captures a parsed Go source package
type PkgSource struct {
	SrcDir  string            // GOPATH/src or GOROOT/src/pkg
	FileSet *token.FileSet    // File names are relative to SrcDir
	PkgPath string
	Pkgs    map[string]*ast.Package
	MainPkg *ast.Package      // Package named after the containing source directory, or main
}

func (p *PkgSource) Name() string {
	_, name := path.Split(p.PkgPath)
	return name
}

// ParsePkg parses package pkg, using FileSet fset
func (l *Layout) ParsePkg(pkgPath string, includeGoRoot bool, mode parser.Mode) (pkgSrc *PkgSource, err error) {
	
	// Find source root for pkgPath
	var srcDir string
	if srcDir, err = l.FindPkg(pkgPath, includeGoRoot); err != nil {
		return nil, err
	}

	// Save current working directory
	var saveDir string
	if saveDir, err = os.Getwd(); err != nil {
		return nil, err
	}

	// Change current directory to root of sources
	if err = os.Chdir(srcDir); err != nil {
		return nil, err
	}
	defer func() {
		err = os.Chdir(saveDir)
	}()

	// Make file set just for this package
	fset := token.NewFileSet()

	// Parse
	var pkgs map[string]*ast.Package
	if pkgs, err = parser.ParseDir(fset, pkgPath, filterGoNoTest, mode); err != nil {
		return nil, err
	}

	// Find the exported package
	var mainPkg *ast.Package
	_, pkgDirName := path.Split(pkgPath)
	for pkgName, pkg := range pkgs {
		if pkgName == pkgDirName || pkgName == "main" {
			mainPkg = pkg
			break
		}
	}
	// TODO: Package source directories will often contain files with main or xxx_test package clauses.
	// We ignore those, by guessing they are not part of the program.
	// The correct way to ignore is to recognize the comment directive: // +build ignore

	return &PkgSource{
		SrcDir:  srcDir,
		FileSet: fset,
		PkgPath: pkgPath,
		Pkgs:    pkgs,
		MainPkg: mainPkg,
	}, nil
}

func filterGoNoTest(fi os.FileInfo) bool {
	n := fi.Name()
	return len(n) > 0 && strings.HasSuffix(n, ".go") && n[0] != '_' && strings.Index(n, "_test.go") < 0
}
