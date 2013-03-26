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

// NewError creates a simple text-based error that is registered with package
// encoding/gob and therefore can be used in places of error interfaces during
// cross-calls. In contrast, note that due to the rules of gob encoding error objects
// that are not explicitly registered with gob cannot be assigned to error interfaces
// that are to be gob-serialized during a cross-call.
func NewError(fmt_ string, arg_ ...interface{}) error {
	return &errorString{fmt.Sprintf(fmt_, arg_...)}
}

// FlattenError converts any error into a gob-serializable one that can be used in cross-calls.
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
