// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proto

import (
	"circuit/kit/xor"
	"encoding/binary"
	"hash/fnv"
)

type tagValue struct {
	TagID
	ValueID
}

type sortTagValues []tagValue

func (stv sortTagValue) Len() int {
	return len(stv)
}

func (stv sortTagValue) Less(i, j int) bool {
	if stv[i].TagID == stv[j].TagID {
		return stv[i].ValueID < stv[j].ValueID
	}
	return stv[i].TagID < stv[j].TagID
}

func (stv sortTagValue) Swap(i, j int) {
	stv[i], stv[j] = stv[j], stv[i]
}

// SpaceID is a unique identifier for the tuple of metric and tags
type SpaceID uint64

func (id SpaceID) ShardKey() xor.Key {
	return xor.Key(id)
}

func HashSpace(m MetricID, t map[TagID]ValueID) SpaceID {
	h := fnv.New64a()
	var tags sortTagValues
	for k, v := range t {
		tags = append(tags, tagValue{k, v})
	}
	sort.Sort(tags)
	for _, tv := range tags {
		??
	}
	return ValueID(h.Sum64())
}

// MetricID is a unique identifier for a metric name
type MetricID uint32

func HashMetric(m string) MetricID {
	h := fnv.New32a()
	h.Write(m)
	return MetricID(h.Sum32())
}

// TagID is the type of integral IDs that string tag key values are hashed to
type TagID uint32

func HashTag(t string) TagID {
	h := fnv.New32a()
	h.Write(t)
	return TagID(h.Sum32())
}

// ValueID is the type of integral IDs that string tag values are hashed to.
// The zero value represents a wildcard tag value in a query context.
type ValueID uint32

func HashValue(t string) ValueID {
	h := fnv.New32a()
	h.Write(t)
	return ValueID(h.Sum32())
}

type XAdd struct {
	MetricID MetricID
	Time     int64
	Tags     map[TagID]ValueID
	Value    float64
}

type XQuery struct {
	MetricID         MetricID
	MinTime, MaxTime int64
	Tags             map[TagID]ValueID
	Stat             Stat			// Statistic SUM or AVG
	Velocity         bool                   // Output first derivative of statistic
}

const (
	Sum Stat = iota
	Avg
)
