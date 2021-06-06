package constant

import "time"

var (
	minTime = time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
)

func MinTime() time.Time {
	return minTime
}

func MaxTime() time.Time {
	return maxTime
}
