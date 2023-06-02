package reddit

import (
	"encoding/json"
	"fmt"
	"gobot/freegames"
	"net/http"
)

// Get Top reddit post from reddit API. It will return array of Post struct
func getPost(subreddit string) (*[]Post, error) {
	resp, err := http.Get("https://www.reddit.com/r/" + subreddit + "/top/.json")
	if err != nil {
		fmt.Println("Failed to get subreddit data:", err)
		return nil, err
	}
	var data Response
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

/*
Filter function for FreeGamesOnSteam subreddit.
Keep only post that fulfill these conditions:

- Created after 7 days ago

- Upvote more than 200

- Do not have ended flair
*/
func freeSteamFilter(data *[]Post) []Post {
	sevenDaysAgo := freegames.DaysUnix(-7)
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

/*
Filter function for FreeGamesFinding subreddit.
Keep only post that fulfill these conditions:

- Created after 7 days ago

- Upvote more than 300

- Do not have any of these flair; Mod post, Regional Issues, Expired
*/
func freeFindingFilter(data *[]Post) []Post {
	sevenDaysAgo := freegames.DaysUnix(-7)
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

/*
GetPostsAndFilter
Get Posts from FreeGamesOnSteam and FreeGameFindings, after that applied filter for each subreddit

- FreeGamesOnSteam 👉freeSteamFilter

- FreeGameFindings 👉freeFindingFilter
*/
func GetPostsAndFilter() []Post {
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
