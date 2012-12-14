package hbase

import (
	"io"
	"testing"
)

type friend struct {
	Follower uint64
	Followee uint64
	Time     int64
}

func TestReader(t *testing.T) {
	r, err := OpenFile("testdata/friends")
	if err != nil {
		t.Fatalf("open (%s)", err)
	}
	var v friend
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
