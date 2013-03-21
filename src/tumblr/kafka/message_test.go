package kafka

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	testMessage = &Message{
		Compression: NoCompression,
		Payload:     []byte{1, 2, 3},
	}
)

func TestMessage(t *testing.T) {
	var w bytes.Buffer
	m0 := testMessage
	if err := m0.Write(&w); err != nil {
		t.Fatalf("message write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	m1 := &Message{}
	_, err := m1.Read(r)
	if err != nil {
		t.Fatalf("message read (%s)", err)
	}
	if !reflect.DeepEqual(m0, m1) {
		t.Errorf("expecting %v got %v", m0, m1)
	}
}
