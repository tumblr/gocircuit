package c

import (
	"os"
	"testing"
	"circuit/c/source"
)

var (
	testLayout = source.NewLayout(os.Getenv("GOROOT"), source.GoPaths{"./testdata"}, "")
)

func TestBuild(t *testing.T) {
	b, err := NewBuild(testLayout, "./_tmp")
	if err != nil {
		t.Fatalf("build (%s)", err)
	}
	if err := b.Build("b"); err != nil {
		t.Errorf("run (%s)", err)
	}
}
