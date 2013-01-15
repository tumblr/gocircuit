package types

import (
	"circuit/c/errors"
	"circuit/c/util"
	"go/ast"
	"go/token"
	"path"
	"reflect"
)

// pkgPath is the package path where the type expression expr resides.
// fimp is the import structure of the source file containing expr.
func CompileTypeSpec(pkgPath string, spec *ast.TypeSpec, fimp *util.FileImports) (typ Type, err error) {
	return compileTypeExpr(pkgPath, spec.Type, fimp)
}

func compileTypeExpr(pkgPath string, expr ast.Expr, fimp *util.FileImports) (typ Type, err error)

	switch q := expr.(type) {

	// Built-in types or references to other types in this package
	case *ast.Ident:
		switch q.Name {
		case "bool":
			typ = Builtin[Bool]
		case "int":
			typ = Builtin[Int]
		case "int8":
			typ = Builtin[Int8]
		case "int16":
			typ = Builtin[Int16]
		case "int32":
			typ = Builtin[Int32]
		case "int64":
			typ = Builtin[Int64]
		case "uint":
			typ = Builtin[Uint]
		case "uint8":
			typ = Builtin[Uint8]
		case "uint16":
			typ = Builtin[Uint16]
		case "uint32":
			typ = Builtin[Uint32]
		case "uint64":
			typ = Builtin[Uint64]
		case "uintptr":
			typ = Builtin[Uintptr]
		case "float32":
			typ = Builtin[Float32]
		case "float64":
			typ = Builtin[Float64]
		case "complex64":
			typ = Builtin[Complex64]
		case "complex128":
			typ = Builtin[Complex128]
		case "string":
			typ = Builtin[String]
		default:
			// Name of another type defined in this package
			typ = &Link{pkgPath, q.Name}
		}
		return typ, nil

	case *ast.ParenExpr:
		return compileTypeExpr(pkgPath, q, fimp)

	case *ast.SelectorExpr:
		pkgAlias, ok := q.X.(*ast.Ident)
		if !ok {
			panic("package alias is not an identifier")
		}
		typeName := q.Sel.Name
		impPath, ok := fimp.Alias[pkgAlias]
		if !ok {
			return nil, errors.New("no import with given alias")
		}
		return &Link{impPath, typeName}, nil

	case *ast.StarExpr:
		base, err := compileTypeExpr(pkgPath, q.X, fimp)
		if err != nil {
			return nil, err
		}
		return &Pointer{Base: base}, nil

	case *ast.ArrayType:
		elt, err := compileTypeExpr(pkgPath, q.Elt, fimp)
		if err != nil {
			return nil, err
		}
		if q.Len == nil {
			return &Slice{Elt: elt}, nil
		}
		if _, ok := q.Len.(*ast.Ellipses); ok {
			return &Array{Len: -1/*XXX*/, Elt: elt}, nil
		}
		return nil, errors.New("unknown array length")

	case *ast.ChanType:
		value, err := compileTypeExpr(pkgPath, q.Value, fimp)
		if err != nil {
			return nil, err
		}
		return &Chan{Dir: q.Dir, Elt: value}, nil

	case *ast.FuncType:
		// TODO: Compile signature details
		return &Signature{}, nil

	case *ast.InterfaceType:
		// TODO: Compile interface details
		return &Interface{}, nil

	case *ast.MapType:
		key, err := compileTypeExpr(pkgPath, q.Key, fimp)
		if err != nil {
			return nil, err
		}
		value, err := compileTypeExpr(pkgPath, q.Value, fimp)
		if err != nil {
			return nil, err
		}
		return &MapType{Key: key, Value: value}, nil

	case *ast.StructType:
		// TODO: Compile struct details
		return &Struct{}, nil
	}

	return nil, errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
}
