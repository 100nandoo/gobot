package freegames

import (
	"gobot/pkg"
)

/*
Merge 2 things:

- Posts from API response

- Posts from supabase

It will keep all the posts from API response that is not inside supabase
*/
func merge() []Post {
	var results []Post
	posts := GetPostsAndFilter()
	supabasePosts := GetAllPost()
	for i := 0; i < len(posts); i++ {
		found := false
		for _, value := range supabasePosts {
			if value.URL == posts[i].URL {
				found = true
				break
			}
		}
		if !found {
			results = append(results, posts[i])
		}
	}
	return results
}

/*
Scouting
Run every day at 11:00
*/
func Scouting() {
	pkg.EverydayAtThisHour(func() {
		pkg.LogWithTimestamp("Scouting free games")
		var merged = merge()
		for _, post := range merged {
			SendPost(post)
			Insert(post)
		}
	}, "11:00")
}
