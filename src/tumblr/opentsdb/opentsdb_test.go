package opentsdb

import (
	"testing"
)

func TestPut(t *testing.T) {
	conn, err := Dial("opentsdb.datacenter.net")
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	err = conn.Put("hello.world", 5, Tag{"host", "test.host"})
	if err != nil {
		t.Errorf("put (%s)", err)
	}
}
