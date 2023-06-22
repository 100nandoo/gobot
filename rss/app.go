package rss

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"gobot/pkg"
	"sync"
)

// Get Rss feed from remote Ok API. It will return array of Item
func getRssItems(rss SupabaseRss) ([]*gofeed.Item, error) {
	feed, err := pkg.Parser.ParseURL(rss.Url)
	if err != nil {
		fmt.Println("Error calling getRssItems", err)
		return nil, err
	}
	return feed.Items, nil
}

/*
Filter function for Rss Feed.
Keep only item not is not older than 7 days ago
*/
func filter(items []*gofeed.Item) *[]gofeed.Item {
	var result []gofeed.Item
	sevenDaysAgo := pkg.DaysUnix(-7)

	for _, item := range items {
		if item.PublishedParsed.Unix() > sevenDaysAgo {
			result = append(result, *item)
		}
	}
	return &result
}

/*
getItemsAndFilter
Get Jobs from remote Ok API, after that applied filter
*/
func getItemsAndFilter(supabaseRss SupabaseRss) *[]gofeed.Item {
	result := make([]gofeed.Item, 0)

	items, err := getRssItems(supabaseRss)

	if err != nil {
		fmt.Println("Error calling getRssItems", err)
		return nil
	}
	result = *filter(items)

	return &result
}

func Scouting() {
	pkg.EverySaturdayDayAtThisHour(func() {
		var wg sync.WaitGroup
		var result []gofeed.Item

		fmt.Println("Scouting Rss")
		rss := GetAllSupabaseRss()

		fmt.Println("rss", rss)
		for _, supabaseRss := range rss {
			wg.Add(1)
			supabaseRss := supabaseRss
			go func() {
				defer wg.Done()
				result = append(result, *getItemsAndFilter(supabaseRss)...)
			}()
		}
		wg.Wait()

		for _, item := range result {
			SendRssItem(item)
		}
	}, "11:05")
}
