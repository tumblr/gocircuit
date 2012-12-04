package c

import (
	"testing"
)

func TestBuild(t *testing.T) {
	b, err := NewBuild("./tmp")
	if err != nil {
		t.Fatalf("build (%s)", err)
	}
	if err := b.Build("circuit/load"); err != nil {
		t.Errorf("run (%s)", err)
	}
}
