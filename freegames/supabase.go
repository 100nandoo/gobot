package freegames

import (
	"fmt"
	"github.com/nedpals/supabase-go"
	"gobot/config"
	"os"
)

type SupabasePost struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Sent  bool   `json:"sent"`
	Dummy bool   `json:"dummy"`
}

var Client = supabase.CreateClient(os.Getenv(config.SupabaseUrl), os.Getenv(config.SupabaseKey))

/*
GetAllPost

Get All rows from Games Database, return arrays of SupabasePost
*/
func GetAllPost() []SupabasePost {
	var results []SupabasePost
	err := Client.DB.From("Games").Select("*").Execute(&results)
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
	err := Client.DB.From("Games").Insert(SupabasePost{
		URL:   post.URL,
		Title: post.Title,
		Sent:  true,
	}).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Insert", err)
		return
	}
}
