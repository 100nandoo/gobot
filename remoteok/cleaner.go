package remoteok

import (
	"gobot/pkg"
	"time"
)

/*
Return list of SupabaseJob that is older than 7 days ago
*/
func findOldSupabaseJobs(supabaseJobs []SupabaseJob) []SupabaseJob {
	var result []SupabaseJob
	sevenDaysAgoUnix := pkg.DaysUnix(-7)
	for _, job := range supabaseJobs {
		if int64(job.Epoch) < sevenDaysAgoUnix {
			result = append(result, job)
		}
	}
	return result
}

/*
Cleaning
Run every saturday at 10:30
1. Find SupabaseJob(s) that is older than 7 days ago
2. Delete it
*/
func Cleaning(now bool) {
	cleaningLogic := func() {
		pkg.LogWithTimestamp("Cleaning remoteOk")
		supabaseJobs := GetAllSupabaseJob()
		oldPosts := findOldSupabaseJobs(supabaseJobs)
		for _, post := range oldPosts {
			Delete(post)
		}
	}

	if now {
		cleaningLogic()
	} else {
		pkg.SpecificDayAtThisHour(cleaningLogic, time.Saturday, "10:30")
	}
}