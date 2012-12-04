package c

import (
	"testing"
)

func TestDep(t *testing.T) {
	l, err := NewWorkingLayout() // This will give us a gopath to the current circuit repo
	if err != nil {
		t.Fatalf("layout (%s)", err)
	}
	dep, err := l.CompileDep("circuit/load")
	if err != nil {
		t.Fatalf("compute dep (%s)", err)
	}
	for _, d := range dep {
		println(d)
	}
}