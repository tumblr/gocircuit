package c

import (
	"go/ast"
)

type Parser interface {
	ParsePkg(pkgPath string) (map[string]*ast.Package, error)
}

// DepTable maintains the dependent packages for a list of incrementally added
// target packages
type DepTable struct {
	parser Parser
	pkgs   map[string]*depPkg
	follow []string
}

type depPkg struct {
	imports  []string
}

func NewDepTable(parser Parser) *DepTable {
	return &DepTable{
		parser: parser,
		pkgs:   make(map[string]*depPkg),
		follow: nil,
	}
}

func (dt *DepTable) Add(pkg string) error {
	dt.follow = append(dt.follow, pkg)
	return dt.loop()
}

func (dt *DepTable) loop() error {
	for len(dt.follow) > 0 {
		pop := dt.follow[0]
		dt.follow = dt.follow[1:]

		// Check if package already processed
		if _, present := dt.pkgs[pop]; present {
			continue
		}

		// Parse package source
		pkgs, err := dt.parser.ParsePkg(pop)
		if err != nil {
			return err
		}

		// Process all import specs in all source files
		imps := make(map[string]struct{})
		for _, pkg := range pkgs {
			pimps := pkgImports(pkg)
			for i, _ := range pimps {
				if i != "C" {
					imps[i] = struct{}{}
				}
			}
		}

		// Make pkg structure and enqueue new imports
		dpkg := &depPkg{}
		for pkg, _ := range imps {
			dpkg.imports = append(dpkg.imports, pkg)
			dt.follow = append(dt.follow, pkg)
		}

		// Save pkg structure
		dt.pkgs[pop] = dpkg
	}
	return nil
}

func (dt *DepTable) All() []string {
	var all []string
	for pkg, _ := range dt.pkgs {
		all = append(all, pkg)
	}
	return all
}
