package main

import (
	"fmt"
	supa "github.com/nedpals/supabase-go"
)

type SupabasePost struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Sent    bool   `json:"sent"`
	FoundAt string `json:"found_at"`
	Dummy   bool   `json:"dummy"`
}

func getAllSupabase() []SupabasePost {
	client := supa.CreateClient(SupabaseUrl, SupabaseKey)

	var results []SupabasePost
	err := client.DB.From("Games").Select("*").Execute(&results)
	if err != nil {
		fmt.Println("Error select all", err)
		return nil
	}
	return results
}

func insert(post Post) {

}
