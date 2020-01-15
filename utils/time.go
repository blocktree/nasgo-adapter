package utils

import (
	"fmt"
	"time"
)

// GetEpochTime return the time span in seconds
func GetEpochTime() int64 {
	d := BeginEpochTime()
	return int64(time.Since(d)) / 1000000
}

func GetTime(t int64) time.Time {
	dur := time.Duration(t * 1000000)
	d := BeginEpochTime()
	return d.Add(dur)
}

func BeginEpochTime() time.Time {
	var d, err = time.Parse("2006-01-02 15:04 MST", "2018-02-04 20:00 UTC")
	if err != nil {
		fmt.Errorf("%s", err)
	}
	return d
}
