package opentsdb

import (
	"testing"
	"time"
)

func TestBestPut(t *testing.T) {
	conn, err := BestEffortDial("opentsdb.datacenter.net")
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	time.Sleep(3 * time.Second)
	err = conn.Put("wiktor.is.in.love", 5, Tag{"host", "test.host"})
	if err != nil {
		t.Errorf("put (%s)", err)
	}
}
