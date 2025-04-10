package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SpotifyTrackResponse struct {
	Album            Album           `json:"album"`
	Artists          []Artist        `json:"artists"`
	AvailableMarkets []string        `json:"available_markets"`
	DiscNumber       int             `json:"disc_number"`
	DurationMs       int             `json:"duration_ms"`
	Explicit         bool            `json:"explicit"`
	ExternalIDs      ExternalIDs     `json:"external_ids"`
	ExternalURLs     ExternalURLs    `json:"external_urls"`
	Href             string          `json:"href"`
	ID               string          `json:"id"`
	IsPlayable       bool            `json:"is_playable"`
	Restrictions     Restrictions    `json:"restrictions"`
	Name             string          `json:"name"`
	Popularity       int             `json:"popularity"`
	PreviewURL       string          `json:"preview_url"`
	TrackNumber      int             `json:"track_number"`
	Type             string          `json:"type"`
	URI              string          `json:"uri"`
	IsLocal          bool            `json:"is_local"`
}

type Album struct {
	AlbumType         string        `json:"album_type"`
	TotalTracks       int           `json:"total_tracks"`
	AvailableMarkets  []string      `json:"available_markets"`
	ExternalURLs      ExternalURLs  `json:"external_urls"`
	Href              string        `json:"href"`
	ID                string        `json:"id"`
	Images            []Image       `json:"images"`
	Name              string        `json:"name"`
	ReleaseDate       string        `json:"release_date"`
	ReleaseDatePrecision string     `json:"release_date_precision"`
	Restrictions      Restrictions  `json:"restrictions"`
	Type              string        `json:"type"`
	URI               string        `json:"uri"`
	Artists           []Artist      `json:"artists"`
}

type Artist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ExternalIDs struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

type Restrictions struct {
	Reason string `json:"reason"`
}


// GetSpotifyTrack sends a GET request to fetch track details from the Spotify API
func GetSpotifyTrack(trackID, accessToken string) (*SpotifyTrackResponse, error) {
	// Prepare the URL
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the provided access token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Spotify API error: %s", errorResp["error"].(map[string]interface{})["message"])
	}

	// Parse the response body
	var trackResp SpotifyTrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&trackResp); err != nil {
		return nil, err
	}

	return &trackResp, nil
}
