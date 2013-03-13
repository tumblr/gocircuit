package api

import (
	"hash/fnv"
	"encoding/json"
	"math"
	"time"
	"circuit/app/sumr"
)

// A Feature is an expressive form of a counter key, which hashes down to the latter
type Feature map[string]string

func (f Feature) Key() sumr.Key {
	buf := []byte(f.String())
	hash := fnv.New64a()
	hash.Write(buf)
	g := hash.Sum(nil)
	var k uint64
	for i := 0; i < 64/8; i++ {
		k |= uint64(g[i]) << uint(i*8)
	}
	return sumr.Key(k)
}

func (f Feature) String() string {
	buf, err := json.Marshal(f)
	if err != nil {
		panic("feature marshal")
	}
	return string(buf)
}

func MakeFeatureMap(b map[string]interface{}) (Feature, error) {
	f := make(Feature)
	for k, v := range b {
		s, ok := v.(string)
		if !ok {
			return nil, ErrFieldType
		}
		f[k] = s
	}
	return f, nil
}

// Change combines a feature together with an event of a given value and time
type Change struct {
	Time    time.Time
	Feature Feature
	Value   float64
}

// Key returns the hash key corresponding to the change's feature
func (s *Change) Key() sumr.Key {
	return s.Feature.Key()
}

// ReadChange parses a change from its JSON representation, like so:
// 
//	{
//		"t": 12345678,
//		"f": { "fkey": "fvalue", ... },
//		"v": 1.234
//	}
//
func ReadChange(dec *json.Decoder) (*Change, error) {
	b := make(map[string]interface{})
	if err := dec.Decode(&b); err != nil {
		return nil, err
	}
	return MakeChangeMap(b)
}

func MakeChangeMap(b map[string]interface{}) (*Change, error) {
	// Read time
	time_, ok := b["t"]
	if !ok {
		return nil, ErrNoValue
	}
	timef, ok := time_.(float64)
	if !ok {
		return nil, ErrNoValue
	}
	if math.IsNaN(timef) || timef < 0 {
		return nil, ErrTime
	}
	t := time.Unix(0, int64(timef))

	// Read value
	value_, ok := b["v"]
	if !ok {
		return nil, ErrNoValue
	}
	value, ok := value_.(float64)
	if !ok {
		return nil, ErrNoValue
	}

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
	return &Change{Time: t, Feature: f, Value: value}, nil
}
