// Copyright 2012 Tumblr, Inc.
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

package kafka

import (
	"time"
)

// Handy time constants for use in Kafka client invokations
const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day

	Earliest = -2
	Latest   = -1
)

// TimeToKafka converts t to the Kafka time format
func TimeToKafka(t time.Time) int64 {
	return t.UnixNano() / 1e9
}

// Now returns the current time in the Kafka format
func Now() int64 {
	return TimeToKafka(time.Now())
}
