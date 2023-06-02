package main

type RedditResponse struct {
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
