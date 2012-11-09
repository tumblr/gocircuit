package block

import (
	"reflect"
	"testing"
	"time"
)

func TestCode(t *testing.T) {
	a := &add{time.Now(), 1, 0.1}
	b := encodeAdd(a)
	println(len(b))
	a2, err := decodeAdd(b)
	if err != nil {
		t.Errorf("decode add (%s)", err)
	}
	if !reflect.DeepEqual(a, a2) {
		t.Errorf("mismatch")
	}
}
