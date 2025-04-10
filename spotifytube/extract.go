package spotifytube

import (
	"fmt"
	"strings"
)
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

func ExtractYoutubeVideoID(url string) (string, error) {
	const param = "v="
	if idx := strings.Index(url, "youtu.be/"); idx != -1 {
		start := idx + len("youtu.be/")
		path := url[start:]
		if parts := strings.SplitN(path, "?", 2); len(parts) > 0 {
			if parts[0] == "" {
				return "", fmt.Errorf("invalid video ID")
			}
			return parts[0], nil
		}
	}
	if idx := strings.Index(url, param); idx != -1 {
		start := idx + len(param)
		path := url[start:]
		if parts := strings.SplitN(path, "&", 2); len(parts) > 0 {
			if parts[0] == "" {
				return "", fmt.Errorf("invalid video ID")
			}
			return parts[0], nil
		}
	}
	return "", fmt.Errorf("unable to extract video ID")
}
