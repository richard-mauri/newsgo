package dashboard

import (
	"context"
	"github.com/richard-mauri/newsgo/internal/feeds"
	"sync/atomic"
	"testing"
	"time"
)

func TestSharedThumbnailFallbackAndCache(t *testing.T) {
	sharedThumbnailCache.Lock()
	sharedThumbnailCache.entries = map[string]sharedThumbnailEntry{}
	sharedThumbnailCache.Unlock()
	entries := []feeds.Article{{Source: "Publisher One", URL: "https://example.com/one"}, {Source: "Publisher One", URL: "https://example.com/two", Image: "https://example.com/feed.jpg"}}
	var hits atomic.Int64
	lookup := func(_ context.Context, address string) string { hits.Add(1); return "https://example.com/image.jpg" }
	enrichSharedImages(context.Background(), entries, lookup)
	if hits.Load() != 1 || entries[0].Image != "https://example.com/image.jpg" || entries[1].Image != "https://example.com/feed.jpg" {
		t.Fatalf("unexpected enrichment: %#v, hits=%d", entries, hits.Load())
	}
	next := []feeds.Article{{Source: "Publisher One", URL: "https://example.com/one"}}
	enrichSharedImages(context.Background(), next, lookup)
	if hits.Load() != 1 || next[0].Image == "" {
		t.Fatal("cached image not reused")
	}
}
func TestSharedThumbnailPerPublisherBudget(t *testing.T) {
	sharedThumbnailCache.Lock()
	sharedThumbnailCache.entries = map[string]sharedThumbnailEntry{}
	sharedThumbnailCache.Unlock()
	a := make([]feeds.Article, 25)
	for i := range a {
		a[i] = feeds.Article{Source: "The Publisher", URL: "https://example.com/a" + string(rune('A'+i))}
	}
	var calls atomic.Int64
	enrichSharedImages(context.Background(), a, func(context.Context, string) string { calls.Add(1); return "" })
	if calls.Load() != 8 {
		t.Fatalf("expected 8 bounded lookups; got %d", calls.Load())
	}
	sharedThumbnailCache.Lock()
	for _, v := range sharedThumbnailCache.entries {
		if time.Since(v.at) > time.Minute {
			t.Fatal("cache timestamp stale")
		}
	}
	sharedThumbnailCache.Unlock()
}
