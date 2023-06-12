package pkg

import (
	"fmt"
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

/*
SupabaseDateToUnix

# Convert supabaseDate to Unix Timestamp

Example:

freegames.SupabaseDateToUnix("2023-06-06") = 1686009600
*/
func SupabaseDateToUnix(dateString string) (int64, error) {
	layout := "2006-01-02"

	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error calling SupabaseDateToUnix", err)
		return 0, err
	}

	unixTimestamp := date.Unix()
	return unixTimestamp, nil
}

/*
NowSupabaseDate

# Return today in supabase date format

Example:

freegames.NowSupabaseDate() = "2023-06-12"
*/
func NowSupabaseDate() string {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	return formattedDate
}
