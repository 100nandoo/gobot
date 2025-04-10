package spotifytube

import (
	"testing"
)

func TestExtractSpotifyTrackID(t *testing.T) {
	tests := []struct {
		url         string
		expectedID  string
		expectError bool
	}{
		{"https://open.spotify.com/track/1ek8UP8J0cHPVx9vGIztSi?si=d5a085a6c2e14300", "1ek8UP8J0cHPVx9vGIztSi", false},
		{"https://open.spotify.com/track/5VJNF9RdCPN99IDCbmMchz", "5VJNF9RdCPN99IDCbmMchz", false},
		{"https://open.spotify.com/track/3N7ZtjvOotNS8AvQraBPrC?si=aabbcc", "3N7ZtjvOotNS8AvQraBPrC", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			id, err := ExtractSpotifyTrackID(test.url)
			if test.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != test.expectedID {
				t.Errorf("expected %s, got %s", test.expectedID, id)
			}
		})
	}
}

func TestExtractYoutubeVideoID(t *testing.T) {
	tests := []struct {
		url         string
		expectedID  string
		expectError bool
	}{
		{"https://www.youtube.com/watch?v=SX_ViT4Ra7k", "SX_ViT4Ra7k", false},
		{"https://music.youtube.com/watch?v=SX_ViT4Ra7k", "SX_ViT4Ra7k", false},
		{"https://www.youtube.com/watch?v=SX_ViT4Ra7k&list=PLNRPV4mrGOshXzInyb5X7DnVMNhOVfeIX", "SX_ViT4Ra7k", false},
		{"https://music.youtube.com/watch?v=8BUIzLo1Dmk&list=PLcgpoLDUNq9hnkN9SFPUNirqmTaXhNje2", "8BUIzLo1Dmk", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			id, err := ExtractYoutubeVideoID(test.url)
			if test.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id != test.expectedID {
				t.Errorf("expected %s, got %s", test.expectedID, id)
			}
		})
	}
}
