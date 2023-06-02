package main

import (
	"gobot/freegames/reddit"
	"gobot/freegames/supabase"
	"gobot/freegames/telegram"
)

/*
Merge 2 things:

- Posts from API response

- Posts from supabase

It will keep all the posts from API response that is not inside supabase
*/
func merge() []reddit.Post {
	var results []reddit.Post
	posts := reddit.GetPostsAndFilter()
	supabasePosts := supabase.GetAllPost()
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

func main() {
	var merged = merge()
	for _, post := range merged {
		telegram.SendPost(post)
		supabase.Insert(post)
	}
}
