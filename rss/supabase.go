package rss

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"gobot/pkg"
)

type SupabaseRss struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Category string `json:"category"`
}

const dbName = "Rss"

/*
GetAllSupabaseRss

Get All rows from Rss Database, return arrays of SupabaseRss
*/
func GetAllSupabaseRss() []SupabaseRss {
	var result []SupabaseRss
	err := pkg.SupabaseClient.DB.From(dbName).Select("*").Execute(&result)
	if err != nil {
		fmt.Println("Error calling GetAllSupabaseRss", err)
		return nil
	}
	return result
}

/*
Insert

Insert a row into Rss Database
*/
func Insert(feed *gofeed.Item) {
	var results []gofeed.Feed
	err := pkg.SupabaseClient.DB.From(dbName).Insert(SupabaseRss{
		Url:      feed.Link,
		Name:     feed.Title,
		Priority: 0,
		Category: "",
	}).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Insert", err)
		return
	}
}

/*
Delete

Delete a row in Rss Database
*/
func Delete(rss SupabaseRss) {
	var results []SupabaseRss
	err := pkg.SupabaseClient.DB.From(dbName).Delete().Eq("url", rss.Url).Execute(&results)
	if err != nil {
		fmt.Println("Error calling Delete", err)
		return
	}
}
