package api

import (
	"circuit/app/sumr"
	"encoding/json"
	"io"
)

// Response is the common response object
type response struct {
	Sum float64 `json:"sum"`
}

// ReadRequestBatchFunc reads a request batch
type readRequestBatchFunc func(io.Reader) ([]interface{}, error)

// Request to add a new change to a feature vector
// On the wire, it looks like so
//
//	{
//		"f": { "fkey": "fvalue", ... },
//		"v": 1.234
//	}
//
type addRequest struct {
	change *change
}

func (r *addRequest) Key() sumr.Key {
	return r.change.Key()
}

func (r *addRequest) Value() float64 {
	return r.change.Value
}

func readAddRequest(dec *json.Decoder) (interface{}, error) {
	change, err := readChange(dec)
	if err != nil {
		return nil, err
	}
	return &addRequest{change: change}, nil
}

func readAddRequestBatch(r io.Reader) ([]interface{}, error) {
	dec := json.NewDecoder(r)
	var bch []interface{}
	for {
		r, err := readAddRequest(dec)
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
type sumRequest struct {
	feature feature
}

func (r *sumRequest) Key() sumr.Key {
	return r.feature.Key()
}

func readSumRequest(dec *json.Decoder) (interface{}, error) {
	b := make(map[string]interface{})
	if err := dec.Decode(&b); err != nil {
		return nil, err
	}
	return makeSumRequestMap(b)
}

func readSumRequestBatch(r io.Reader) ([]interface{}, error) {
	dec := json.NewDecoder(r)
	var bch []interface{}
	for {
		r, err := readSumRequest(dec)
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

func makeSumRequestMap(b map[string]interface{}) (*sumRequest, error) {
	// Read feature
	feature_, ok := b["f"]
	if !ok {
		return nil, ErrNoFeature
	}
	feature, ok := feature_.(map[string]interface{})
	if !ok {
		return nil, ErrNoFeature
	}
	f, err := makeFeatureMap(feature)
	if err != nil {
		return nil, err
	}

	// Done
	return &sumRequest{feature: f}, nil
}
