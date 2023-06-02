package main

import (
	"fmt"
	"gobot/freegames/reddit"
	"gobot/freegames/supabase"
)

/*
Merge 2 things:

- Posts from API response

- Posts from supabase

It will keep all the posts from API response that is not inside supabase
*/
func merge() []reddit.Post {
	posts := reddit.GetPostsAndFilter()
	supabasePosts := supabase.GetAllPost()
	for i := 0; i < len(posts); i++ {
		for _, value := range supabasePosts {
			if value.URL == posts[i].URL {
				posts = append(posts[:i], posts[i+1:]...)
				break
			}
		}
	}
	return posts
}

// Insert arrays of reddit.Post to Supabase
func insertToSupabase(mergedPosts []reddit.Post) {
	for _, post := range mergedPosts {
		supabase.Insert(post)
	}
}

func main() {
	insertToSupabase(merge())
	supabasePosts := supabase.GetAllPost()
	for _, post := range supabasePosts {
		fmt.Println(post)
	}
}
