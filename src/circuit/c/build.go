package c

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

type Build struct {
	layout  *Layout
	jail    *Jail

	deps    []string
}

func NewBuild(jaildir string) (b *Build, err error) {

	b = &Build{}

	// Derive the user source layout from the environment
	if b.layout, err = NewWorkingLayout(); err != nil {
		return nil, err
	}

	// Create a new compilation jail
	if b.jail, err = NewJail(jaildir); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Build) Build(pkgs ...string) error {

	var err error

	// Calculate package dependencies
	if b.deps, err = b.layout.CompileDep(pkgs...); err != nil {
		return err
	}

	// Process each package
	for _, pkg := range b.deps {

		// Parse package
		skel, err := b.layout.ParsePkg(pkg, parser.ParseComments)
		if err != nil {
			return err
		}

		// Print each file in the package in turn
		for _, pkgAST := range skel.Pkgs {
			for fileName, fileAST := range pkgAST.Files {
				if err := b.printFile(pkg, fileName, skel.FileSet, fileAST); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (b *Build) printFile(pkg, fileName string, fset *token.FileSet, file *ast.File) error {
	f, err := b.jail.CreateSrcFile(pkg, fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return printer.Fprint(f, fset, file)
}
