package circuit

import (
	"encoding/gob"
	"errors"
	"fmt"
)

func init() {
	gob.Register(&errorString{})
	gob.Register(errors.New(""))
}

// NewError creates a simple text-based error that is serializable
func NewError(fmt_ string, arg_ ...interface{}) error {
	return &errorString{fmt.Sprintf(fmt_, arg_...)}
}

func FlattenError(err error) error {
	if err == nil {
		return nil
	}
	return NewError(err.Error())
}

type errorString struct {
	S string
}

func (e *errorString) Error() string {
	return e.S
}
