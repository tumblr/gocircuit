package hbase

import (
	"io"
	"testing"
)

type record struct {
	Field1 uint64
	Field2 uint64
	Field3 int64
}

func TestReader(t *testing.T) {
	r, err := OpenFile("testdata/records")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	var v record
	for {
		err = r.Read(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read (%s)", err)
		}
		println(v.Follower, v.Followee, v.Time)
	}
}
