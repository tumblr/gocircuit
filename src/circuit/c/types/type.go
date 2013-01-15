package types

import (
	"go/ast"
	"go/token"
	"math/big"
)

// Type is a type definition.
type Type interface {
	aType()
}

// Incomplete type specializations

// Link is an unresolved type reference
type Link struct {
	PkgPath string
	Name    string
}

func (*Link) aType() {}

// TypeExpr stands in for an unparsed type expression
type TypeExpr struct{
	expr *ast.TypeExpr
}

type TypeSource struct {
	FileSet *token.FileSet
	Spec    *ast.TypeSpec
	PkgPath string
}

// Type specializations

type Array struct {
	Len int64
	Elt Type
}

// Basic is a type definition
type Basic struct {
	Kind BasicKind
	Info BasicInfo
	Size int64
	Name string
}

var Builtin = [...]*Basic{
	Invalid:        {aType, Invalid,        0,                      0,  "invalid type"},

	Bool:           {aType, Bool,           IsBoolean,              1,  "bool"},
	Int:            {aType, Int,            IsInteger,              0,  "int"},
	Int8:           {aType, Int8,           IsInteger,              1,  "int8"},
	Int16:          {aType, Int16,          IsInteger,              2,  "int16"},
	Int32:          {aType, Int32,          IsInteger,              4,  "int32"},
	Int64:          {aType, Int64,          IsInteger,              8,  "int64"},
	Uint:           {aType, Uint,           IsInteger | IsUnsigned, 0,  "uint"},
	Uint8:          {aType, Uint8,          IsInteger | IsUnsigned, 1,  "uint8"},
	Uint16:         {aType, Uint16,         IsInteger | IsUnsigned, 2,  "uint16"},
	Uint32:         {aType, Uint32,         IsInteger | IsUnsigned, 4,  "uint32"},
	Uint64:         {aType, Uint64,         IsInteger | IsUnsigned, 8,  "uint64"},
	Uintptr:        {aType, Uintptr,        IsInteger | IsUnsigned, 0,  "uintptr"},
	Float32:        {aType, Float32,        IsFloat,                4,  "float32"},
	Float64:        {aType, Float64,        IsFloat,                8,  "float64"},
	Complex64:      {aType, Complex64,      IsComplex,              8,  "complex64"},
	Complex128:     {aType, Complex128,     IsComplex,              16, "complex128"},
	String:         {aType, String,         IsString,               0,  "string"},
	UnsafePointer:  {aType, UnsafePointer,  0,                      0,  "Pointer"},

	UntypedBool:    {aType, UntypedBool,    IsBoolean | IsUntyped,  0,  "untyped boolean"},
	UntypedInt:     {aType, UntypedInt,     IsInteger | IsUntyped,  0,  "untyped integer"},
	UntypedRune:    {aType, UntypedRune,    IsInteger | IsUntyped,  0,  "untyped rune"},
	UntypedFloat:   {aType, UntypedFloat,   IsFloat   | IsUntyped,  0,  "untyped float"},
	UntypedComplex: {aType, UntypedComplex, IsComplex | IsUntyped,  0,  "untyped complex"},
	UntypedString:  {aType, UntypedString,  IsString  | IsUntyped,  0,  "untyped string"},
	UntypedNil:     {aType, UntypedNil,     IsUntyped,              0,  "untyped nil"},
}

// BasicInfo stores auxiliary information about a basic type
type BasicInfo int

const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)

// BasicKind distinguishes a primitive type
type BasicKind int

const (
	Invalid BasicKind = iota

	// Predeclared types
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	// Types for untyped values
	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	// Aliases
	Byte = Uint8
	Rune = Int32
)

type Chan struct {
	Dir ast.ChanDir
	Elt Type
}

type Field struct {
	Name        string
	Type        Type
	Tag         string
	IsAnonymous bool
}

type Interface struct {
	Methods []*Method
}

type Map struct {
	Key, Value Type
}

type Method struct{
	Name string
	Type *Signature
}

type Named struct {
	Name       string
	PkgPath    string
	Underlying Type
}

func (n *Named) FullName() string {
	return n.PkgPath + "·" + n.Name
}

type Nil struct{}

type Pointer struct {
	Base Type
}

type Result struct {
	Values []Type
}

type Signature struct {
	Recv       Type
	Params     []Type
	Results    []Type
	IsVariadic bool
}

type Slice struct {
	Elt Type
}

type Struct struct {
	Fields []*Field
}

type ComplexConstant struct {
	Re, Im *big.Rat
}
