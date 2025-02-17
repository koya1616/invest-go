package utils

import "time"

func IsWithinTimeRange(now time.Time, start int, end int) bool {
	currentMinutes := now.Hour()*60 + now.Minute()
	return currentMinutes >= start && currentMinutes <= end
}
