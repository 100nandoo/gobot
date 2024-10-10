package freegames

import (
	"encoding/json"
	"gobot/pkg"
	"net/http"
)

type Response struct {
	Data struct {
		Children []struct {
			Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Post struct {
	ApprovedAtUtc interface{} `json:"approved_at_utc"`
	Subreddit     string      `json:"subreddit"`
	Title         string      `json:"title"`
	Name          string      `json:"name"`
	UpvoteRatio   float64     `json:"upvote_ratio"`
	Ups           int         `json:"ups"`
	LinkFlairText string      `json:"link_flair_text"`
	Score         int         `json:"score"`
	Thumbnail     string      `json:"thumbnail"`
	Created       float64     `json:"created"`
	ID            string      `json:"id"`
	Author        string      `json:"author"`
	URL           string      `json:"url"`
	CreatedUtc    float64     `json:"created_utc"`
	Media         interface{} `json:"media"`
	IsVideo       bool        `json:"is_video"`
}

// Get Top reddit post from reddit API. It will return array of Post struct
func getPost(subreddit string) (*[]Post, error) {
	resp, err := http.Get("https://www.reddit.com/r/" + subreddit + "/top/.json")
	if err != nil {
		pkg.LogWithTimestamp("Failed to get subreddit data: %v", err)
		return nil, err
	}
	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		pkg.LogWithTimestamp("Failed to decode JSON: %v", err)
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
	sevenDaysAgo := pkg.DaysUnix(-7)
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
	sevenDaysAgo := pkg.DaysUnix(-7)
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
		pkg.LogWithTimestampInt(value.Score, value.Title, value.URL)
	}

	return result
}
