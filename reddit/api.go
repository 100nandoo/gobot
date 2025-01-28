package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// RedditResponse represents the structure of Reddit's JSON response
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// Post represents a Reddit post
type Post struct {
	Title         string        `json:"title"`
	Image         string        `json:"url_overridden_by_dest"`
	IsGallery     bool          `json:"is_gallery"`
	Score         int           `json:"score"`
	Author        string        `json:"author"`
	Permalink     string        `json:"permalink"`
	Created       float64       `json:"created"`
	MediaMetadata MediaMetadata `json:"media_metadata"`
	Images        []string      // New field to store gallery image URLs
}

// MediaMetadata represents the media_metadata field in a Reddit post
type MediaMetadata map[string]struct {
	M string `json:"m"` // MIME type, e.g., "image/jpeg"
	S struct {
		U string `json:"u"` // URL of the image
	} `json:"s"`
	ID string `json:"id"`
}

// TimeFilter represents the time period for top posts
type TimeFilter string

const (
	Month TimeFilter = "month"
	Week  TimeFilter = "week"
)

// FetchTopPosts fetches top posts from a specified subreddit for a given time period and filters posts with score above 100
func FetchTopPosts(subreddit string, timeFilter TimeFilter, score int) (*RedditResponse, error) {
	url := fmt.Sprintf("https://old.reddit.com/r/%s/top.json?t=%s", subreddit, timeFilter)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var redditResp RedditResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Filter posts with score above 100
	filteredChildren := []struct {
		Data Post `json:"data"`
	}{}
	for _, child := range redditResp.Data.Children {
		if child.Data.Score > score {
			filteredChildren = append(filteredChildren, child)
		}
	}

	// Update the response with filtered posts
	redditResp.Data.Children = filteredChildren

	return &redditResp, nil
}

// ConvertMediaMetadataToURLs converts MediaMetadata to an array of URLs
func ConvertMediaMetadataToURLs(metadata MediaMetadata) []string {
	var urls []string

	// Iterate over each entry in the MediaMetadata
	for _, data := range metadata {
		// Extract the file extension from the "m" field (e.g., "image/jpg" -> "jpg")
		fileExt := strings.Split(data.M, "/")
		if len(fileExt) != 2 {
			continue // Skip malformed MIME types
		}

		// Construct the URL in the format https://i.redd.it/IDm
		url := fmt.Sprintf("https://i.redd.it/%s.%s", data.ID, fileExt[1])
		urls = append(urls, url)
	}

	return urls
}

// UnmarshalJSON implements a custom unmarshal function for Post
func (p *Post) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursive UnmarshalJSON calls
	type PostAlias Post
	alias := &struct {
		*PostAlias
	}{
		PostAlias: (*PostAlias)(p),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Populate Images field based on media_metadata or image URL
	if p.IsGallery {
		p.Images = ConvertMediaMetadataToURLs(p.MediaMetadata)
	} else if p.Image != "" {
		p.Images = []string{p.Image}
	}

	return nil
}
