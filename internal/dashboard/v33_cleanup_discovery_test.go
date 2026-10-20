package dashboard

import (
	"github.com/richard-mauri/newsgo/internal/feeds"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV33FeedSearchPublisherRanking(t *testing.T) {
	if feedMatchRank("New Scientist", "new scientist") != 0 {
		t.Fatal("exact publisher search must lead")
	}
	if feedMatchRank("New Scientist Weekly", "new scientist") >= feedMatchRank("News about Science", "new scientist") {
		t.Fatal("publisher prefix must rank before unrelated result")
	}
}

func TestV33NoCachedCatalogHeadlinesOrServerRefresh(t *testing.T) {
	app, err := New([]feeds.Source{{Name: "Example", Category: "News", URL: "https://example.org/rss"}}, "https://scrapergo.example", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	html := w.Body.String()
	if strings.Contains(html, "Fetching headlines…") || strings.Contains(html, "cached headlines") {
		t.Fatal("obsolete server-side headline status in dashboard")
	}
	if !strings.Contains(html, "Waiting for subscribed feeds") {
		t.Fatal("missing client-side status")
	}
	if !strings.Contains(html, "grid-template-columns:minmax(0,1fr) auto") {
		t.Fatal("missing layout safeguard")
	}
}

func TestV33NewScientistURLSafeForDiscovery(t *testing.T) {
	if !safeURL("https://www.newscientist.com/feed/") {
		t.Fatal("publisher's canonical feed URL rejected")
	}
}
