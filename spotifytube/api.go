package spotifytube

import (
	"encoding/json"
	"fmt"
	"gobot/config"
	"io"
	"net/http"
	"os"
	"strings"
)

type InspectResponse struct {
	Status string `json:"status"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Data   struct {
		ExternalID  string   `json:"externalId"`
		PreviewURL  *string  `json:"previewUrl"`
		Name        string   `json:"name"`
		ArtistNames []string `json:"artistNames"`
		AlbumName   string   `json:"albumName"`
		ImageURL    string   `json:"imageUrl"`
		ISRC        *string  `json:"isrc"`
		Duration    int      `json:"duration"`
		URL         string   `json:"url"`
	} `json:"data"`
}

type SearchResponse struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Source string     `json:"source"`
	Status string     `json:"status"`
	Data   SearchData `json:"data"`
	Type   string     `json:"type"`
}

type SearchData struct {
	ExternalID  string   `json:"externalId"`
	PreviewURL  *string  `json:"previewUrl"`
	Name        string   `json:"name"`
	ArtistNames []string `json:"artistNames"`
	AlbumName   string   `json:"albumName"`
	ImageURL    string   `json:"imageUrl"`
	ISRC        *string  `json:"isrc"`
	Duration    int      `json:"duration"`
	URL         *string  `json:"url"`
}

const (
	baseUrl        = "https://api.musicapi.com/public/"
	inspectUrlPath = "inspect/url"
	searchUrlPath  = "search"

	POST                 = "POST"
	YOUTUBE_MUSIC_PREFIX = "https://music.youtube.com/watch?v="
)

// makeAPIRequest is a reusable function to handle API requests
func makeAPIRequest(endpoint, payload string) ([]byte, error) {
	client := &http.Client{}
	url := baseUrl + endpoint

	req, err := http.NewRequest(POST, url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add necessary headers
	token := "Token " + os.Getenv(config.MusicAPItoken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", token)

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}

// Search performs a track search on the API and prints the parsed response
// Define a custom type for source selection
type SourceType int

// Define constants for the possible source options
const (
	SPOTIFY_ONLY SourceType = iota
	YOUTUBE_ONLY
	BOTH
)

func Search(track string, artist string, sourceType SourceType) (*SearchResponse, error) {
	var sources string

	// Determine the sources based on the sourceType
	switch sourceType {
	case SPOTIFY_ONLY:
		sources = `["spotify"]`
	case YOUTUBE_ONLY:
		sources = `["youtubeMusic"]`
	case BOTH:
		sources = `["youtubeMusic", "spotify"]`
	default:
		return nil, fmt.Errorf("invalid source type")
	}

	// Build the payload with the selected sources
	payload := fmt.Sprintf(`{
		"track": "%s",
		"artist": "%s",
		"type": "track",
		"sources": %s
	}`, track, artist, sources)

	// Make the API request
	body, err := makeAPIRequest(searchUrlPath, payload)
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}

	// Parse the response into SearchResponse struct
	var searchResp SearchResponse
	err = json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// formattedResp, err := json.MarshalIndent(searchResp, "", "  ")
	// if err != nil {
	// 	log.Println("Error marshalling JSON:", err)
	// 	return nil, err
	// }
	// fmt.Println(string(formattedResp))
	return &searchResp, nil
}

// InspectUrl retrieves information about a specific track by URL and prints the parsed response
func InspectUrl(url string) (*InspectResponse, error) {
	// Build the payload
	payload := fmt.Sprintf(`{
		"url": "%s"
	}`, url)

	// Make the API request
	body, err := makeAPIRequest(inspectUrlPath, payload)
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}

	// Parse the response into InspectResponse struct
	var inspectResp InspectResponse
	err = json.Unmarshal(body, &inspectResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Return the parsed InspectResponse and nil error
	return &inspectResp, nil
}
