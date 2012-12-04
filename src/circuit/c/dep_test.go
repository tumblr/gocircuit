package c

import (
	"testing"
)

func TestDep(t *testing.T) {
	build, err := NewWorkingBuild() // This will give us a gopath to the current circuit repo
	if err != nil {
		t.Fatalf("build (%s)", err)
	}
	dep, err := CompileDep(gopath, "circuit/load")
	if err != nil {
		t.Fatalf("compute dep (%s)", err)
	}
	for _, d := range dep {
		println(d)
	}
}
