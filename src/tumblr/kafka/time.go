package kafka

import (
	"time"
)

// Handy time constants for use in Kafka client invokations
const (
	Second = 1
	Minute = 60*Second
	Hour   = 60*Minute
	Day    = 24*Hour
	Week   = 7*Day

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
