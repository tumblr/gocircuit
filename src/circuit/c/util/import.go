package util

import (
	"go/ast"
	"path"
	"strconv"
)

// DeterminePkgImports returns a map of all package paths, directly imported by pkg
func DeterminePkgImports(pkg *ast.Package) map[string]struct{} {
	imprts := make(map[string]struct{}) 
	for _, file := range pkg.Files {
		for _, impSpec := range file.Imports {
			_, importPath := parseImportSpec(impSpec)
			imprts[importPath] = struct{}{}
		}
	}
	return imprts
}

// DetermineFileImports …
func DetermineFileImports(file *ast.File) (alias map[string]string, dot, underscore []string) {
	for _, impSpec := range file.Imports {
		pkgAlias, pkgPath := parseImportSpec(impSpec)
		switch pkgAlias {
		case ".":
			dot = append(dot, pkgPath)
		case "_":
			underscore = append(underscore, pkgPath)
		case "":
			panic("import with no alias")
		default:
			alias[pkgAlias] = pkgPath
		}
	}
	return
}

func parseImportSpec(spec *ast.ImportSpec) (pkgAlias, pkgPath string) {
	var err error
	if pkgPath, err = strconv.Unquote(spec.Path.Value); err != nil {
		panic(err)
	}
	if spec.Name == nil {
		_, pkgAlias = path.Split(pkgPath)
	} else {
		pkgAlias = spec.Name.Name
	}
	return
}
