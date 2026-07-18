package reddit

import "testing"

func TestConvertMediaMetadataToURLsPreservesGalleryOrder(t *testing.T) {
	gallery := GalleryData{
		Items: []GalleryItem{
			{MediaID: "second"},
			{MediaID: "first"},
		},
	}
	metadata := MediaMetadata{
		"first": {
			M: "image/jpg",
			S: struct {
				U string `json:"u"`
			}{
				U: "https://preview.redd.it/first.jpg?width=100&amp;format=pjpg",
			},
			ID: "first",
		},
		"second": {
			M: "image/png",
			S: struct {
				U string `json:"u"`
			}{
				U: "https://preview.redd.it/second.png?width=100&amp;format=png",
			},
			ID: "second",
		},
	}

	got := ConvertMediaMetadataToURLs(gallery, metadata)

	want := []string{
		"https://preview.redd.it/second.png?width=100&format=png",
		"https://preview.redd.it/first.jpg?width=100&format=pjpg",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d urls, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got url[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestConvertMediaMetadataToURLsFallsBackToConstructedURL(t *testing.T) {
	metadata := MediaMetadata{
		"image": {
			M:  "image/jpeg",
			ID: "abc123",
		},
	}

	got := ConvertMediaMetadataToURLs(GalleryData{}, metadata)

	if len(got) != 1 {
		t.Fatalf("got %d urls, want 1", len(got))
	}
	if got[0] != "https://i.redd.it/abc123.jpeg" {
		t.Fatalf("got %q, want fallback i.redd.it URL", got[0])
	}
}
