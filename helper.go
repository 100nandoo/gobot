package main

import "time"

func daysUnix(days int) int64 {
	now := time.Now()

	// Subtract 7 days from the current time
	sevenDaysAgo := now.AddDate(0, 0, days)

	// Get the Unix timestamp for 7 days ago
	unixTimestamp := sevenDaysAgo.Unix()

	return unixTimestamp
}
