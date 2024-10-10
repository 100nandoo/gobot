package remoteok

import (
	"gobot/pkg"
)

/*
Merge 2 things:

- Jobs from API response

- Jobs from supabase

It will keep all the jobs from API response that is not inside supabase
*/
func merge() []Job {
	var results []Job
	jobs := getJobsAndFilter()
	supabaseJobs := GetAllSupabaseJob()
	for i := 0; i < len(jobs); i++ {
		found := false
		for _, value := range supabaseJobs {
			if value.URL == jobs[i].URL {
				found = true
				break
			}
		}
		if !found {
			results = append(results, jobs[i])
		}
	}
	return results
}

/*
Scouting
Run every day at 10:00
*/
func Scouting() {
	pkg.EverydayAtThisHour(func() {
		pkg.LogWithTimestamp("Scouting remoteOk")
		var merged = merge()
		for _, job := range merged {
			SendJob(job)
			Insert(job)
		}
	}, "10:00")
}
