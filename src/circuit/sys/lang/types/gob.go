package types

import (
	"encoding/gob"
	"reflect"
	"time"
)

// Register some common types. Repeated registration is ok.
func init() {
	gob.Register(make(map[string]interface{}))
	gob.Register(make(map[string]string))
	gob.Register(make(map[string]int))
	gob.Register(make([]interface{}, 0))
	gob.Register(time.Duration(0))
}

// gobFlattenRegister registers the flattened type of t with gob
// E.g. the flattened type of *T is T, of **T is T, etc.
// Interface types cannot be registered.
func gobFlattenRegister(t reflect.Type) {
	if t.Kind() == reflect.Interface {
		return
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pz := reflect.New(t)
	gob.Register(pz.Elem().Interface())
}
