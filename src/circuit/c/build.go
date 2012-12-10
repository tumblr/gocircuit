package c

import (
	"circuit/c/errors"
	"circuit/c/types"
	"go/ast"
	"go/parser"
)

type Build struct {
	layout    *Layout
	jail      *Jail

	pkgs      map[string]*PkgSrc  // pkgPath to parsed package source
	depTable  *DepTable
	typeTable *types.TypeTable
}

func NewBuild(layout *Layout, jaildir string) (b *Build, err error) {

	b = &Build{layout: layout}

	// Create a new compilation jail
	if b.jail, err = NewJail(jaildir); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Build) Build(pkgPaths ...string) error {

	var err error

	// Calculate package dependencies
	b.pkgs = make(map[string]*PkgSrc)
	if err = b.compileDep(pkgPaths...); err != nil {
		return err
	}

	// Parse types
	b.typeTable = types.NewTypeTable()
	if err = b.parseTypes(); err != nil {
		return err
	}

	// dbg
	for _, typ := range b.typeTable.ListFullNames() {
		println(typ)
	}

	if err = b.transformRegisterInterfaces(); err != nil {
		return err
	}

	// Flush rewritten source into jail
	if err = b.flush(); err != nil {
		return err
	}

	return nil
}

// ParsePkg parses the requested package path and saves the resulting package
// AST node into the pkgs field
func (b *Build) ParsePkg(pkgPath string) (map[string]*ast.Package, error) {
	pkgSrc, err := b.layout.ParsePkg(pkgPath, false, parser.ParseComments)
	if err != nil {
		Log("- %s skipping", pkgPath)
		// This is intended for Go's packages itself, which we don't want to parse for now
		return nil, nil
	}
	Log("+ %s", pkgPath)

	// Save package AST into global map
	if _, present := b.pkgs[pkgPath]; present {
		return nil, errors.New("package %s already parsed", pkgPath)
	}
	b.pkgs[pkgPath] = pkgSrc

	return pkgSrc.Pkgs, nil
}

// compileDep causes all packages that pkgs depend on to be parsed
func (b *Build) compileDep(pkgPaths ...string) error {
	Log("Calculating dependencies ...")
	Indent()
	defer Unindent()

	b.depTable = NewDepTable(b)
	for _, pkgPath := range pkgPaths {
		if err := b.depTable.Add(pkgPath); err != nil {
			return err
		}
	}
	return nil
}

// parseTypes finds all type declarations and registers them with a global map
func (b *Build) parseTypes() error {
	for pkgPath, pkgSrc := range b.pkgs {
		if pkgSrc.MainPkg == nil {
			continue
		}
		if err := b.typeTable.AddPkg(pkgSrc.FileSet, pkgPath, pkgSrc.MainPkg); err != nil {
			return err
		}
	}
	return nil
}
