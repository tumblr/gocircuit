package api

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadBatch(t *testing.T) {
	const src = `{"f":{"a":"b"},"v":1}{"f":{"c":"d"}, "v":-1}`
	r := bytes.NewBufferString(src)
	req, err := ReadAddRequestBatch(Now(), r)
	if err != nil {
		t.Errorf("read batch (%s)", err)
	}
	if len(req) != 2 {
		t.Errorf("wrong number of read requests")
	}
	for _, r := range req {
		fmt.Printf("%T\n", r)
	}
}
