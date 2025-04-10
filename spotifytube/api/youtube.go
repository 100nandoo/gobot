package api

import (
	"context"
	"gobot/config"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func InitYoutube() (*youtube.Service, error){
	token :=  os.Getenv(config.YoutubeToken)
	
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(token))
	
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	return service, err
}

func SearchYoutube(service youtube.Service, query string) (*youtube.SearchListResponse, error) {
	
	call := service.Search.List([]string{"id"}).
		Q(query).
		MaxResults(3).
		Type("video")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error searching for %q: %v", query, err)
	}
	return response, err
}