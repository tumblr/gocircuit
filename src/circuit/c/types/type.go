package types

import (
	"circuit/c/errors"
	"circuit/c/util"
	"go/ast"
	"go/token"
	"path"
	"reflect"
)

type Type struct {

	// Sweep 1
	FileSet *token.FileSet
	Spec    *ast.TypeSpec
	Name    string
	PkgPath string

	// Sweep 2
	Kind    reflect.Kind      // Go kind
	Elem    *Type             // Ptr, Interface, Array, Slice
	Procs   []*Proc           // Methods
}

func (t *Type) FullName() string {
	return path.Join(t.PkgPath, t.Name)
}

// UnlinkedType represents a partially-compiled type.
// If Kind equals:
//	reflect.Invalid, the type is an alias for another type, whose name resides in Elem.
//	reflect.Ptr, reflect.Slice, reflect.Array, the name of the element type resides in Elem.
//	reflect.Map, the name of the key (value) type is stored in Key (Value).
type UnlinkedType struct {
	Kind  reflect.Kind
	Elem  string	// Fully-qualified name of another type
	Key   string	// "
	Value string	// "
}

func compileTypeSpec(spec *ast.TypeSpec) (unlinked *Type0, err error) {
	…
	compileTypeExpr(spec.Type)
}

func compileTypeExpr(pkgPath string, expr ast.Expr, fimp *util.FileImports) (unlinked *UnlinkedType, err error)
	unilnked = &UnlinkedType{}
	switch q := expr.(type) {

	// Built-in types or references to other types in this package
	case *ast.Ident:
		switch q.Name {
		case "bool":
			unlinked.Kind = reflect.Bool
		case "int":
			unlinked.Kind = reflect.Int
		case "int8":
			unlinked.Kind = reflect.Int8
		case "int16":
			unlinked.Kind = reflect.Int16
		case "int32":
			unlinked.Kind = reflect.Int32
		case "int64":
			unlinked.Kind = reflect.Int64
		case "uint":
			unlinked.Kind = reflect.Uint
		case "uint8":
			unlinked.Kind = reflect.Uint8
		case "uint16":
			unlinked.Kind = reflect.Uint16
		case "uint32":
			unlinked.Kind = reflect.Uint32
		case "uint64":
			unlinked.Kind = reflect.Uint64
		case "uintptr":
			unlinked.Kind = reflect.Uintptr
		case "float32":
			unlinked.Kind = reflect.Float32
		case "float64":
			unlinked.Kind = reflect.Float64
		case "complex64":
			unlinked.Kind = reflect.Complex64
		case "complex128":
			unlinked.Kind = reflect.Complex128
		case "string":
			unlinked.Kind = reflect.String
		default:
			// Name of another type defined in this package
			unlinked.AliasFor = path.Join(pkgPath, q.Name)
		}
		return unlinked, nil

	case *ast.ParenExpr:
		return compileTypeExpr(pkgPath, q, fimp)

	case *ast.SelectorExpr:
		pkgAlias, ok := q.X.(*ast.Ident)
		if !ok {
			panic("unrecognized selector")
		}
		typeName := q.Sel.Name
		impPath, ok := fimp.Alias[pkgAlias]
		if !ok {
			return nil, errors.New("import alias unknown")
		}
		unlinked.AliasFor = path.Join(impPath, typeName)
		return unlinked, nil

	case *ast.StarExpr:
		// r.Elem will be filled in during a follow up sweep of all types
		unlinked.Kind = reflect.Ptr
		?

	case *ast.ArrayType:
		XX // Slice or array kind?
		unlinked.Kind = reflect.Array

	case *ast.ChanType:
		unlinked.Kind = reflect.Chan

	case *ast.FuncType:
		unlinked.Kind = reflect.Func

	case *ast.InterfaceType:
		unlinked.Kind = reflect.Interface

	case *ast.MapType:
		unlinked.Kind = reflect.Map

	case *ast.StructType:
		unlinked.Kind = reflect.Struct

	default:
		return 0, "", errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return unlinked, nil
}
