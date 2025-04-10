package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobot/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SpotifyTokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	RequestedAt time.Time `json:"requested_at"`
}

func GetSpotifyAccessToken() (*SpotifyTokenResponse, error) {
	clientID := os.Getenv(config.SpotifyID)
	clientSecret := os.Getenv(config.SpotifySecret)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	now := time.Now()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get token: %s", string(body))
	}

	var tokenResp SpotifyTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}
	tokenResp.RequestedAt = now

	return &tokenResp, nil
}

func ExtractSpotifyTrackID(url string) (string, error) {
	const prefix = "/track/"
	idx := strings.Index(url, prefix)
	if idx == -1 {
		return "", fmt.Errorf("invalid track URL")
	}
	start := idx + len(prefix)
	path := url[start:]
	if parts := strings.SplitN(path, "?", 2); len(parts) > 0 {
		if parts[0] == "" {
			return "", fmt.Errorf("invalid track ID")
		}
		return parts[0], nil
	}
	return "", fmt.Errorf("unable to extract track ID")
}
