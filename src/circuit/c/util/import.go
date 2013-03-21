package util

import (
	"go/ast"
	"path"
	"strconv"
)

// CompilePkgImports returns a map of all package paths, directly imported by pkg
func CompilePkgImports(pkg *ast.Package) map[string]struct{} {
	imprts := make(map[string]struct{})
	for _, file := range pkg.Files {
		for _, impSpec := range file.Imports {
			_, importPath := parseImportSpec(impSpec)
			imprts[importPath] = struct{}{}
		}
	}
	return imprts
}

type FileImports struct {
	Alias      map[string]string // Import alias to package path
	Dot        []string          // List of package paths imported with the dot alias
	Underscore []string          // List of package paths imported with the underscore alias
}

// CompileFileImports …
func CompileFileImports(file *ast.File) (fimp *FileImports) {
	fimp = &FileImports{Alias: make(map[string]string)}
	for _, impSpec := range file.Imports {
		pkgAlias, pkgPath := parseImportSpec(impSpec)
		switch pkgAlias {
		case ".":
			fimp.Dot = append(fimp.Dot, pkgPath)
		case "_":
			fimp.Underscore = append(fimp.Underscore, pkgPath)
		case "":
			panic("import with no alias")
		default:
			fimp.Alias[pkgAlias] = pkgPath
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
