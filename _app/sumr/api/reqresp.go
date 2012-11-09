package api

import (
	"encoding/json"
	"io"
	"tumblr/app/sumr"
)

// Response is the common response object
type Response struct {
	Sum float64  `json:"sum"`
}

// ReadRequestBatchFunc reads a request batch
type ReadRequestBatchFunc func(io.Reader) ([]interface{}, error)

// Request to add a new change to a feature vector
// On the wire, it looks like so
//
//	{
//		"f": { "fkey": "fvalue", ... },
//		"v": 1.234
//	}
//
type AddRequest struct {
	Change *Change
}

func (r *AddRequest) Key() sumr.Key {
	return r.Change.Key()
}

func (r *AddRequest) Value() float64 {
	return r.Change.Value
}

func ReadAddRequest(dec *json.Decoder) (interface{}, error) {
	change, err := ReadChange(dec)
	if err != nil {
		return nil, err
	}
	return &AddRequest{Change: change}, nil
}

func ReadAddRequestBatch(r io.Reader) ([]interface{}, error) {
	dec := json.NewDecoder(r)
	var bch []interface{}
	for {
		r, err := ReadAddRequest(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bch = append(bch, r)
	}
	return bch, nil
}

// A SumRequest returns the sum of all changes at a given feature
// On the wire, it looks like so
//
//	{
//		"f": { "fkey": "fvalue", ... },
//	}
//
type SumRequest struct {
	Feature Feature
}

func (r *SumRequest) Key() sumr.Key {
	return r.Feature.Key()
}

func ReadSumRequest(dec *json.Decoder) (interface{}, error) {
	b := make(map[string]interface{})
	if err := dec.Decode(&b); err != nil {
		return nil, err
	}
	return MakeSumRequestMap(b)
}

func ReadSumRequestBatch(r io.Reader) ([]interface{}, error) {
	dec := json.NewDecoder(r)
	var bch []interface{}
	for {
		r, err := ReadSumRequest(dec)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bch = append(bch, r)
	}
	return bch, nil
}

func MakeSumRequestMap(b map[string]interface{}) (*SumRequest, error) {
	// Read feature
	feature_, ok := b["f"]
	if !ok {
		return nil, ErrNoFeature
	}
	feature, ok := feature_.(map[string]interface{})
	if !ok {
		return nil, ErrNoFeature
	}
	f, err := MakeFeatureMap(feature)
	if err != nil {
		return nil, err
	}

	// Done
	return &SumRequest{Feature: f}, nil
}
