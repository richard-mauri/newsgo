package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticlePageImage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><meta property="og:image" content="https://img.example.com/story.jpg"></head></html>`))
	}))
	defer srv.Close()
	if got := articlePageImage(context.Background(), srv.Client(), srv.URL); got != "https://img.example.com/story.jpg" {
		t.Fatalf("got %q", got)
	}
}
func TestArticlePageImageAbsent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>No photograph</title></head></html>`))
	}))
	defer srv.Close()
	if got := articlePageImage(context.Background(), srv.Client(), srv.URL); got != "" {
		t.Fatalf("unexpected %q", got)
	}
}
