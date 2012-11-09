package block

import (
	"bytes"
	"encoding/binary"
	"time"
	"tumblr/app/sumr"
)

type add struct {
	UTime time.Time
	Key   sumr.Key
	Value float64
}

type addOnDisk struct {
	UTime  int64
	Key    sumr.Key
	Value  float64
}

// OPTIMIZE: Use a code object that uses the same underlying gob coder on each file

func encodeAdd(a *add) []byte {
	var w bytes.Buffer
	if err := binary.Write(&w, binary.LittleEndian, &addOnDisk{a.UTime.UnixNano(), a.Key, a.Value}); err != nil {
		panic("sumr coder")
	}
	return w.Bytes()
}

func decodeAdd(p []byte) (*add, error) {
	r := bytes.NewBuffer(p)
	_a :=&addOnDisk{}
	if err := binary.Read(r, binary.LittleEndian, _a); err != nil {
		return nil, err
	}
	return &add{time.Unix(0, _a.UTime), _a.Key, _a.Value}, nil
}
