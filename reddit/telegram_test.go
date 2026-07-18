package reddit

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPhotoFromURLDownloadsImageBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte("image-bytes"))
	}))
	defer server.Close()

	photo, err := photoFromURL(server.URL)
	if err != nil {
		t.Fatalf("photoFromURL returned error: %v", err)
	}
	if photo.File.FileURL != "" {
		t.Fatalf("expected uploaded file, got URL-backed file %q", photo.File.FileURL)
	}
	if photo.File.FileReader == nil {
		t.Fatal("expected FileReader to be populated")
	}

	body, err := io.ReadAll(photo.File.FileReader)
	if err != nil {
		t.Fatalf("failed to read photo bytes: %v", err)
	}
	if string(body) != "image-bytes" {
		t.Fatalf("got %q, want %q", string(body), "image-bytes")
	}
}
