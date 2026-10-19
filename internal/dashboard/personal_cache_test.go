package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

func TestPersonalCacheTTL(t *testing.T) {
	if got := personalCacheTTL(cacheEntry{}); got != 10*time.Minute {
		t.Fatalf("success TTL: %s", got)
	}
	if got := personalCacheTTL(cacheEntry{err: "temporary error"}); got != 20*time.Second {
		t.Fatalf("failure TTL: %s", got)
	}
}

func TestCancelledPersonalRequestNeverCaches(t *testing.T) {
	url := "https://example.org/newsgo-cancel-test.xml"
	personalCache.Lock()
	delete(personalCache.entries, url)
	personalCache.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, msg := personalFetch(ctx, feeds.Source{Name: "Test", URL: url})
	if msg != "Request canceled" {
		t.Fatalf("error = %q", msg)
	}
	personalCache.Lock()
	_, exists := personalCache.entries[url]
	personalCache.Unlock()
	if exists {
		t.Fatal("canceled request was cached")
	}
}

func TestCanceledContextDoesNotReturnCachedData(t *testing.T) {
	url := "https://example.org/newsgo-cancel-cached.xml"
	personalCache.Lock()
	personalCache.entries[url] = cacheEntry{at: time.Now()}
	personalCache.Unlock()
	defer func() { personalCache.Lock(); delete(personalCache.entries, url); personalCache.Unlock() }()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, msg := personalFetch(ctx, feeds.Source{Name: "Test", URL: url})
	if msg != "Request canceled" {
		t.Fatalf("error = %q", msg)
	}
}
