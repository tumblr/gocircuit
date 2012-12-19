package c

import (
	"circuit/c/dep"
	"circuit/c/source"
	"circuit/c/types"
	"go/ast"
	"go/parser"
)

type Build struct {
	src   *source.Source
	dep   *dep.Dep
	types *types.TypeTable
}

func NewBuild(layout *source.Layout, writeDir string) (b *Build, err error) {
	src, err := source.New(layout, writeDir)
	if err != nil {
		return nil, err
	}
	return &Build{src: src}, nil
}

func (b *Build) Build(pkgPaths ...string) error {

	var err error

	// Calculate dependencies
	if err = b.determineDep(pkgPaths...); err != nil {
		return err
	}

	// Parse types
	b.types = types.New()
	if err = b.linkTypes(); err != nil {
		return err
	}

	// dbg
	for _, typ := range b.types.ListFullNames() {
		println(typ)
	}

	// Add code that registers all user structs with the circuit runtime type system
	if err = b.TransformRegisterValues(); err != nil {
		return err
	}

	// Flush rewritten source into output jail
	if err = b.src.Flush(); err != nil {
		return err
	}

	return nil
}

type buildParser Build

// Parse implements dep.Parser; It is invoked by the dependency calculator's
// internal algorithm.
func (b *buildParser) Parse(pkgPath string) (map[string]*ast.Package, error) {
	_, inGoRoot, err := b.src.FindPkg(pkgPath)
	if err != nil {
		return nil, err
	}

	// Go packages are not parsed and consecuently their dependencies are not followed
	if inGoRoot {
		return nil, nil
	}

	pkg, _, err := b.src.ParsePkg(pkgPath, parser.ParseComments)
	if err != nil {
		Log("- %s skipping (%s)", pkgPath, err)
		// This is intended for Go's packages itself, which we don't want to parse for now
		return nil, err
	}
	Log("+ %s parsed", pkgPath)

	return pkg.PkgAST, nil
}

// determineDep causes all packages that pkgPaths depend on to be parsed
func (b *Build) determineDep(pkgPaths ...string) error {
	Log("Calculating dependencies ...")
	Indent()
	defer Unindent()

	b.dep = dep.New((*buildParser)(b))
	for _, pkgPath := range pkgPaths {
		if err := b.dep.Add(pkgPath); err != nil {
			return err
		}
	}
	return nil
}

// linkTypes finds all type declarations and registers them with a global map
func (b *Build) linkTypes() error {
	Log("Linking types ...")
	Indent()
	defer Unindent()

	for pkgPath, pkg := range b.src.GetAll() {
		libPkg := pkg.LibPkg()
		if libPkg == nil {
			// XXX: This is probably a main pkg; we still need to
			// link all its types in the worker binary
			continue
		}
		if err := b.types.AddPkg(pkg.FileSet, pkgPath, libPkg); err != nil {
			return err
		}
	}
	return nil
}
