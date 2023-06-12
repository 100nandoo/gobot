package main

import (
	"fmt"
	"gobot/freegames"
)

/*
Merge 2 things:

- Posts from API response

- Posts from supabase

It will keep all the posts from API response that is not inside supabase
*/
func merge() []freegames.Post {
	var results []freegames.Post
	posts := freegames.GetPostsAndFilter()
	supabasePosts := freegames.GetAllPost()
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
	fmt.Println("Gobot started...")
	freegames.Every11am(func() {
		var merged = merge()
		for _, post := range merged {
			freegames.SendPost(post)
			freegames.Insert(post)
		}
	})
}
