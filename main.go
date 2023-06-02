package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getPost(subreddit string) (*[]Post, error) {
	resp, err := http.Get("https://www.reddit.com/r/" + subreddit + "/top/.json")
	if err != nil {
		fmt.Println("Failed to get subreddit data:", err)
		return nil, err
	}
	var data RedditResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Failed to decode JSON:", err)
		return nil, err
	}

	var result []Post
	for _, value := range data.Data.Children {
		result = append(result, value.Post)
	}

	return &result, err
}

func freeSteamFilter(data *[]Post) []Post {
	sevenDaysAgo := daysUnix(-7)
	var result []Post
	for i := 0; i < len(*data); i++ {
		value := (*data)[i]
		if value.Score > 200 &&
			int64(value.Created) > sevenDaysAgo &&
			value.LinkFlairText != "Ended" {
			result = append(result, value)
		}
	}
	return result
}

func freeFindingFilter(data *[]Post) []Post {
	sevenDaysAgo := daysUnix(-7)
	var result []Post
	for i := 0; i < len(*data); i++ {
		value := (*data)[i]
		if value.Score > 300 &&
			int64(value.Created) > sevenDaysAgo &&
			value.LinkFlairText != "Mod Post" &&
			value.LinkFlairText != "Regional Issues" &&
			value.LinkFlairText != "Expired" {
			result = append(result, value)
		}
	}
	return result
}

func getPostsAndFilter() []Post {
	steamPosts, _ := getPost("FreeGamesOnSteam")

	filteredSteamPosts := freeSteamFilter(steamPosts)

	findingPosts, _ := getPost("FreeGameFindings")
	filteredFindingPosts := freeFindingFilter(findingPosts)

	result := append(filteredSteamPosts, filteredFindingPosts...)
	for _, value := range result {
		fmt.Println(value.Score, value.Title, value.URL)
	}

	return result
}

func main() {
	//posts := getPostsAndFilter()
	//supabasePosts := getAllSupabase()
	//for i := 0; i < len(posts); i++ {
	//	for _, value := range supabasePosts {
	//		if value.URL == posts[i].URL {
	//			posts = append(posts[:i], posts[i+1:]...)
	//			break
	//		}
	//	}
	//}
	//fmt.Println(posts)
	Init()
	getAllSupabase()
}
