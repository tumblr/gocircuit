package types


/*
func (tt *TypeTable) linkType(…) … {
	…
	switch q := spec.Type.(type) {
	// Built-in types or references to other types in this package
	case *ast.Ident:
		?
	case *ast.ParenExpr:
		?
	case *ast.SelectorExpr:
		?
	case *ast.StarExpr:
		// r.Elem will be filled in during a follow up sweep of all types
		r.Kind = reflect.Ptr
	case *ast.ArrayType:
		XX // Slice or array kind?
		r.Kind = reflect.Array
	case *ast.ChanType:
		r.Kind = reflect.Chan
	case *ast.FuncType:
		r.Kind = reflect.Func
	case *ast.InterfaceType:
		r.Kind = reflect.Interface
	case *ast.MapType:
		r.Kind = reflect.Map
	case *ast.StructType:
		r.Kind = reflect.Struct
	default:
		return nil, errors.NewSource(fset, spec.Name.NamePos, "unexpected type definition")
	}

	return r
}
*/
