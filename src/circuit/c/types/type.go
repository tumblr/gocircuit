package types

import (
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

type Type0 struct {
	Kind reflect.Kind
	X    string
	Y    string
}

func compileTypeSpec(spec *ast.TypeSpec) (type0 *Type0, err error) {
	…
	compileTypeExpr(spec.Type)
}

func compileTypeExpr(expr ast.Expr) (type0 *Type0, err error)
	type0 = &Type0{}
	switch q := expr.(type) {

	// Built-in types or references to other types in this package
	case *ast.Ident:
		switch q.Name {
		case "bool":
			type0.Kind = reflect.Bool
		case "int":
			type0.Kind = reflect.Int
		case "int8":
			type0.Kind = reflect.Int8
		case "int16":
			type0.Kind = reflect.Int16
		case "int32":
			type0.Kind = reflect.Int32
		case "int64":
			type0.Kind = reflect.Int64
		case "uint":
			type0.Kind = reflect.Uint
		case "uint8":
			type0.Kind = reflect.Uint8
		case "uint16":
			type0.Kind = reflect.Uint16
		case "uint32":
			type0.Kind = reflect.Uint32
		case "uint64":
			type0.Kind = reflect.Uint64
		case "uintptr":
			type0.Kind = reflect.Uintptr
		case "float32":
			type0.Kind = reflect.Float32
		case "float64":
			type0.Kind = reflect.Float64
		case "complex64":
			type0.Kind = reflect.Complex64
		case "complex128":
			type0.Kind = reflect.Complex128
		case "string":
			type0.Kind = reflect.String
		default:
			// Name of another type defined in this package
			type0.Kind = reflect.Invalid
			type0.X = q.Name
		}
		return type0, nil

	case *ast.ParenExpr:
		return compileTypeExpr(q)

	case *ast.SelectorExpr:
		…

	case *ast.StarExpr:
		// r.Elem will be filled in during a follow up sweep of all types
		type0.Kind = reflect.Ptr
		?

	case *ast.ArrayType:
		XX // Slice or array kind?
		type0.Kind = reflect.Array

	case *ast.ChanType:
		type0.Kind = reflect.Chan

	case *ast.FuncType:
		type0.Kind = reflect.Func

	case *ast.InterfaceType:
		type0.Kind = reflect.Interface

	case *ast.MapType:
		type0.Kind = reflect.Map

	case *ast.StructType:
		type0.Kind = reflect.Struct

	default:
		return 0, "", errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return type0, nil
}
