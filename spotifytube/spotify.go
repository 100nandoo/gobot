package spotifytube

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	regexId    = `(?<=track\/)[^&|$|?]*`
	spotifyUrl = "spotify.com"
	queryLimit = 3
)

func urlToSong(url string, client *spotify.Client) *spotify.FullTrack {
	if strings.Contains(url, "spotify.com") {
		id := findID(url)
		fullTrack, err := client.GetTrack(context.Background(), spotify.ID(id))
		if err != nil {
			log.Fatalf("couldn't get track: %v", err)
			return nil
		}
		fmt.Println(fullTrack.Album.Name, fullTrack.Name)
		return fullTrack
	}
	return nil
}

const spotifyIDRegex = `spotify\.com\/(\w+)\/(\w+)`

func findID(url string) string {
	clean := regexp.MustCompile(spotifyIDRegex)
	result := clean.FindStringSubmatch(url)

	if len(result) > 0 {
		return result[2]
	}
	return ""
}

func initClient() *spotify.Client {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client
}
