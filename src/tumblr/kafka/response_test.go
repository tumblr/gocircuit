package kafka

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	testResponseHeader = &ResponseHeader{
		_NonHeaderLen: 123,
		Err:           KafkaErrUnknown,
	}
)

func TestResponseHeader(t *testing.T) {
	var w bytes.Buffer
	rh0 := testResponseHeader
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("response header write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &ResponseHeader{}
	_, err := rh1.Read(r)
	if err != nil {
		t.Errorf("response header read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %v got %v", rh0, rh1)
	}
}

var (
	testFetchResponse = &FetchResponse{
		ResponseHeader: ResponseHeader{
			Err: KafkaErrUnknown,
		},
		Messages: []*Message{
			testMessage,
			testMessage,
		},
	}
)

func TestFetchResponse(t *testing.T) {
	var w bytes.Buffer
	rh0 := testFetchResponse
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("fetch response write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &FetchResponse{}
	_, err := rh1.Read(r)
	if err != nil {
		t.Errorf("fetch response read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %v got %v", rh0, rh1)
	}
}

var (
	testOffsetsResponse = &OffsetsResponse{
		ResponseHeader: ResponseHeader{
			Err: KafkaErrUnknown,
		},
		Offsets: []Offset{
			0x00112233445566,
			0x0f1f2f3f4f5f6f,
		},
	}
)

func TestOffsetsResponse(t *testing.T) {
	var w bytes.Buffer
	rh0 := testOffsetsResponse
	if err := rh0.Write(&w); err != nil {
		t.Fatalf("offsets response write (%s)", err)
	}
	r := bytes.NewBuffer(w.Bytes())
	rh1 := &OffsetsResponse{}
	err := rh1.Read(r)
	if err != nil {
		t.Errorf("offsets response read (%s)", err)
	}
	if !reflect.DeepEqual(rh0, rh1) {
		t.Errorf("expecting %#v got %#v", rh0, rh1)
	}
}
