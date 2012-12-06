package types

import (
	"fmt"
	"go/token"
)

type Error struct {
	FileSet *token.FileSet
	Pos     token.Pos
	Msg     string
}

func NewError(fset *ast.FileSet, pos token.Pos, msg string) *Error {
	return &Error{
		FileSet: set,
		Pos:     pos,
		Msg:     msg,
	}
}

func (e *Error) Error() string {
	pos := e.FileSet.Position(e.Pos)
	return fmt.Sprintf("%s • %s", pos, msg)
}
