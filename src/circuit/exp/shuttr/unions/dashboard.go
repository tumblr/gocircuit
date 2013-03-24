package dashboard

import (
	"bytes"
	"encoding/binary"
	"circuit/kit/xor"
	"tumblr/balkan/proto"
)

// RowKey represents the row keys in the dashboard table
type RowKey struct {
	// It is crucial that TimelineID comes before PostID.
	// This affects the way in which LevelDB keys are serialized.
	TimelineID int64
	PostID     int64
}

func DecodeRowKey(raw []byte) (*RowKey, error) {
	rowKey := &RowKey{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowKey); err != nil {
		return nil, err
	}
	rowKey.PostID *= -1
	return rowKey, nil
}

func (rowKey *RowKey) ShardKey() xor.Key {
	return proto.ShardKeyOf(rowKey.TimelineID)
}

func (rowKey *RowKey) Encode() []byte {
	var w bytes.Buffer
	sortKey := *rowKey
	sortKey.PostID *= -1	// Flipping the sign results in flipping the LevelDB key order
	if err := binary.Write(&w, binary.BigEndian, sortKey); err != nil {
		panic("leveldb dashboard row key encoding")
	}
	return w.Bytes()
}

// RowValue represents the row values in the dashboard table
type RowValue struct {
	PrevPostID int64
}

func DecodeRowValue(raw []byte) (*RowValue, error) {
	rowValue := &RowValue{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowValue); err != nil {
		return nil, err
	}
	return rowValue, nil
}

func (rowValue *RowValue) Encode() []byte {
	var w bytes.Buffer
	if err := binary.Write(&w, binary.BigEndian, rowValue); err != nil {
		panic("leveldb dashboard row value encoding")
	}
	return w.Bytes()
}
