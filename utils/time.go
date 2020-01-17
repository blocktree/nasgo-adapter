package utils

import (
	"time"
)

// GetEpochTime return the time span in seconds
func GetEpochTime() int64 {
	//d := beginEpochTime()
	time := time.Now()
	//return time.Unix() - d.Unix()
	return time.Unix() - 1520193600
}

func beginEpochTime() time.Time {
	d := time.Date(2018, 2, 4, 20, 0, 0, 0, time.UTC)
	return d
}
