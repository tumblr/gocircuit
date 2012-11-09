package lang

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&errorString{})
}

var ErrParse = NewError("parse")

// NewError creates a simple text-based error that is serializable
func NewError(fmt_ string, arg_ ...interface{}) error {
	return &errorString{fmt.Sprintf(fmt_, arg_...)}
}

type errorString struct {
	S string
}

func (e *errorString) Error() string {
	return e.S
}
