package utils

import (
	"time"
)

// GetEpochTime return the time span in seconds
func GetEpochTime() int64 {
	d := BeginEpochTime()
	return int64(time.Since(d)) / 1000000000
}

func GetTime(t int64) time.Time {
	dur := time.Duration(t * 1000000000)
	d := BeginEpochTime()
	return d.Add(dur)
}

func BeginEpochTime() time.Time {
	return time.Date(2018, time.February, 4, 20, 0, 0, 0, time.UTC)
}
