package lang

import (
	"fmt"
	"testing"
)

func TestExportImport(t *testing.T) {
	r := New(NewSandbox())
	x := r.Export(1, 2)
	fmt.Printf("x=%#v\n", x)
	v, s, err := r.Import(x)
	if err != nil {
		t.Errorf("import (%s)", err)
	}
	fmt.Printf("v=%#v, #s=%d\n", v, len(s))
}
