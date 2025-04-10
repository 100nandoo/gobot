package freegames

import (
	"gobot/pkg"
	"time"
)

/*
Return list of SupabasePost that has fountAt older than 7 days ago
*/
func findOldSupabasePosts(supabasePosts []SupabasePost) []SupabasePost {
	var result []SupabasePost
	sevenDaysAgoUnix := pkg.DaysUnix(-7)
	for _, value := range supabasePosts {
		unix, err := pkg.SupabaseDateToUnix(value.FoundAt)
		if err == nil {
			if unix < sevenDaysAgoUnix {
				result = append(result, value)
			}
		}
	}
	return result
}

/*
Cleaning
Run every saturday at 01:00
1. Find SupabasePost(s) that is older than 7 days ago
2. Delete it
*/
func Cleaning(now bool) {
	cleaningLogic := func() {
		pkg.LogWithTimestamp("Cleaning free games")
		supabasePosts := GetAllPost()
		oldPosts := findOldSupabasePosts(supabasePosts)
		for _, post := range oldPosts {
			Delete(post)
		}
	}

	if now {
		cleaningLogic()
	} else {
		pkg.SpecificDayAtThisHour(cleaningLogic, time.Saturday, 1, 0)
	}
}
