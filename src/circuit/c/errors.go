package c

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

func Error(fmt_ string, arg_ ...interface{}) error {
	return errors.New(fmt.Sprintf(fmt_, arg_...))
}
