package pkg

import (
	"time"
)

/*
DaysUnix

# Return n days from now in unix timestamp format

Example:

freegames.DaysUnix(7) = 1686274720, return unix timestamp 7 days from now
*/
func DaysUnix(days int) int64 {
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, days)
	unixTimestamp := sevenDaysAgo.Unix()
	return unixTimestamp
}
