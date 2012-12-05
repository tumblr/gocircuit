package c

import (
	"os"
	"testing"
)

var (
	testLayout = NewLayout(os.Getenv("GOROOT"), GoPaths{"./testdata"}, "")
)

func TestBuild(t *testing.T) {
	b, err := NewBuild(testLayout, "./tmp")
	if err != nil {
		t.Fatalf("build (%s)", err)
	}
	if err := b.Build("b"); err != nil {
		t.Errorf("run (%s)", err)
	}
}
