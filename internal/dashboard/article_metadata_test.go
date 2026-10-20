package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
 "time"
 "github.com/richard-mauri/newsgo/internal/feeds"
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

func TestArticlePageImageBodyFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<html><body><img src="https://cdn.example.com/logo.svg" alt="Logo"><article><img src="https://www.seahorsemagazine.com/images/stories/2026-November/Ronstan_Nov_2026_heading.jpg" alt="Ronstan Nov 2026 heading"></article></body></html>`))
	}))
	defer srv.Close()
	want := "https://www.seahorsemagazine.com/images/stories/2026-November/Ronstan_Nov_2026_heading.jpg"
	if got := articlePageImage(context.Background(), srv.Client(), srv.URL); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
func TestArticlePageImageLazyAndRelative(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<main><img src="data:image/gif;base64,abc" data-src="/photos/sailing.jpg"></main>`))
	}))
	defer srv.Close()
	// HTTP-only test server may fetch the page, but outbound image URLs must be HTTPS.
	if got := articlePageImage(context.Background(), srv.Client(), srv.URL); got != "" {
		t.Fatalf("unexpected insecure image %q", got)
	}
}

func TestArticlePageImageStripsNavigationFragment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<article><img src="https://www.seahorsemagazine.com/images/stories/2026-November/Ronstan_Nov_2026_heading.jpg"></article>`))
	}))
	defer srv.Close()
	got := articlePageImage(context.Background(), srv.Client(), srv.URL+"#new_tab")
	if got != "https://www.seahorsemagazine.com/images/stories/2026-November/Ronstan_Nov_2026_heading.jpg" {
		t.Fatalf("fragment prevented image extraction: %q", got)
	}
}

func TestSeahorseFinnArticleImage(t *testing.T) {
	// ScraperGo extracts this lead image from the same article. NewsGo must
	// discover it without relying on an og:image element or same-site RSS host.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><body><header><img src="https://example.org/logo.svg"></header><main><h1>A rare combination of speed and easy handling</h1><img alt="Finn Nov 2026 heading" src="https://www.seahorsemagazine.com/images/stories/2026-November/Finn_Nov_2026_heading.jpg"></main></body></html>`))
	}))
	defer srv.Close()
	want := "https://www.seahorsemagazine.com/images/stories/2026-November/Finn_Nov_2026_heading.jpg"
	if got := articlePageImage(context.Background(), srv.Client(), srv.URL+"#new_tab"); got != want {
		t.Fatalf("Finn lead-image fallback got %q, want %q", got, want)
	}
}
func TestMissingThumbnailCacheRetriesSooner(t *testing.T) {
	missing := cacheEntry{articles: []feeds.Article{{URL: "https://example.com/story"}}}
	if ttl := personalCacheTTL(missing); ttl > 2*time.Minute {t.Fatalf("image-less success cached too long: %s", ttl)}
}
