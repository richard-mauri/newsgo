package dashboard

import (
	"net/url"
	"testing"
)

func TestAdvertisedFeeds(t *testing.T) {
	base, _ := url.Parse("https://example.com/sports/football")
	html := []byte(`<html><head><link rel="alternate" type="application/rss+xml" title="Football" href="/sports/feed.xml"><link href='/atom.xml' type='application/atom+xml' rel='alternate'><link rel="stylesheet" type="text/css" href="/fake"></head></html>`)
	feeds := advertisedFeeds(html, base)
	if len(feeds) != 2 || feeds[0] != "https://example.com/sports/feed.xml" || feeds[1] != "https://example.com/atom.xml" {
		t.Fatalf("unexpected discovery: %v", feeds)
	}
}
func TestAdvertisedRejectsUnsafeLinks(t *testing.T) {
	base, _ := url.Parse("https://example.com/")
	html := []byte(`<link rel="alternate" type="application/rss+xml" href="http://127.0.0.1/private"><link rel="alternate" type="application/rss+xml" href="javascript:alert(1)">`)
	// safeURL checks scheme; local IP destinations are blocked by the network dialer.
	feeds := advertisedFeeds(html, base)
	if len(feeds) != 1 || feeds[0] != "http://127.0.0.1/private" {
		t.Fatalf("unexpected links: %v", feeds)
	}
}

func TestLinkedFeedsFromDirectory(t *testing.T) {
	base, _ := url.Parse("https://example.org/rss-feeds")
	html := []byte(`<a href="/college-football">College Football</a><a href="/api/optimized-rss?tags=cfb&amp;size=30">College Football RSS</a><a href="https://example.org/atom.xml">News</a>`)
	got := linkedFeeds(html, base)
	if len(got) != 2 || got[0] != "https://example.org/api/optimized-rss?tags=cfb&size=30" || got[1] != "https://example.org/atom.xml" {
		t.Fatalf("unexpected feeds: %v", got)
	}
}
