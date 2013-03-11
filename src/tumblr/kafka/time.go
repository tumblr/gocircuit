package kafka

import (
	"time"
)

const (
	Second = 1
	Minute = 60*Second
	Hour   = 60*Minute
	Day    = 24*Hour
	Week   = 7*Day

	Earliest = -2
	Latest   = -1
)

func TimeToKafka(t time.Time) int64 {
	return t.UnixNano() / 1e9
}

func Now() int64 {
	return TimeToKafka(time.Now())
}
