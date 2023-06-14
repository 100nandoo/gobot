package freegames

import (
	"fmt"
	"gobot/pkg"
)

type SupabasePost struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Sent    bool   `json:"sent"`
	Dummy   bool   `json:"dummy"`
	FoundAt string `json:"found_at"`
}

const dbName = "Games"

/*
GetAllPost

Get All rows from Games Database, return arrays of SupabasePost
*/
func GetAllPost() []SupabasePost {
	var results []SupabasePost
	err := pkg.SupabaseClient.DB.From(dbName).Select("*").Execute(&results)
	if err != nil {
		fmt.Println("Error calling GetAllPost", err)
		return nil
	}
	return results
}

/*
Insert

Insert a row into Games Database
*/
func Insert(post Post) {
	var results []SupabasePost
	err := pkg.SupabaseClient.DB.From(dbName).Insert(SupabasePost{
		URL:     post.URL,
		Title:   post.Title,
		Sent:    true,
		FoundAt: pkg.NowSupabaseDate(),
	}).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Insert", err)
		return
	}
}

/*
Delete

Delete a row in Games Database
*/
func Delete(post SupabasePost) {
	var results []SupabasePost
	err := pkg.SupabaseClient.DB.From(dbName).Delete().Eq("url", post.URL).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Delete", err)
		return
	}
}
