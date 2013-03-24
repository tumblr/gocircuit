package series

import (
	"bytes"
	"circuit/exp/shuttr/proto"
	"circuit/kit/xor"
	"encoding/binary"
)

// Key represents the row key used for the timelines tables in LevelDB
type RowKey struct {
	// It is crucial that TimelineID comes before PostID.
	// This affects the way in which LevelDB keys are serialized.
	TimelineID int64 // TumblelogID of the timeline that is posting
	PostID     int64 // PostID of the new post
}

func DecodeRowKey(raw []byte) (*RowKey, error) {
	rowKey := &RowKey{}
	if err := binary.Read(bytes.NewBuffer(raw), binary.BigEndian, rowKey); err != nil {
		return nil, err
	}
	rowKey.PostID *= -1
	return rowKey, nil
}

// ShardKey returns an xor.Key which determines in which timeline shard this row belongs.
func (rowKey *RowKey) ShardKey() xor.Key {
	return proto.ShardKeyOf(rowKey.TimelineID)
}

// Encode returns the raw LevelDB representation of this row key
func (rowKey *RowKey) Encode() []byte {
	var w bytes.Buffer
	sortKey := *rowKey
	sortKey.PostID *= -1 // Flipping the sign results in flipping the LevelDB key order
	if err := binary.Write(&w, binary.BigEndian, sortKey); err != nil {
		panic("leveldb timeline row key encoding")
	}
	return w.Bytes()
}
