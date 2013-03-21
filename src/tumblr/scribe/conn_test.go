package scribe

import (
	"testing"
)

func TestConn(t *testing.T) {
	conn, err := Dial("devbox:1464")
	if err != nil {
		t.Fatalf("dial (%s)", err)
	}
	if err = conn.Emit([]Message{Message{"test-cat", "test-msg"}}...); err != nil {
		t.Errorf("emit (%s)", err)
	}
	if err = conn.Close(); err != nil {
		t.Errorf("close (%s)", err)
	}
}
