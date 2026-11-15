package dashboard

import (
	"strings"
	"testing"

	"github.com/richard-mauri/newsgo/internal/feeds"
)

func TestArsCanonicalFeedIdentity(t *testing.T) {
	body := []byte(`<?xml version="1.0"?><rss xmlns:atom="http://www.w3.org/2005/Atom"><channel><title>Ars Technica</title><atom:link href="https://arstechnica.com/feed/" rel="self" type="application/rss+xml" /></channel></rss>`)
	if got := feedIdentity(body); got != "https://arstechnica.com/feed" {
		t.Fatalf("canonical feed identity = %q", got)
	}
}

func TestFeedIdentityRequiresSameContentNotSamePublisherTitle(t *testing.T) {
	a, b, other := []feeds.Article{}, []feeds.Article{}, []feeds.Article{}
	for _, slug := range []string{"a", "b", "c", "d", "e", "f"} {
		article := feeds.Article{Title: "Ars Technica", URL: "https://arstechnica.com/story/" + slug}
		a = append(a, article)
		b = append(b, article)
		other = append(other, feeds.Article{Title: article.Title, URL: "https://arstechnica.com/science/" + slug})
	}
	if !equivalentFeeds(a, b, "", "") { t.Fatal("matching article URLs should be equivalent") }
	if equivalentFeeds(a, other, "", "") { t.Fatal("publisher name alone must not merge distinct feeds") }
	if !equivalentFeeds(a, other, "https://arstechnica.com/feed", "https://arstechnica.com/feed") { t.Fatal("matching publisher-declared identity should be equivalent") }
}

func TestFeedIdentityAndAccountDiagnosticsUI(t *testing.T) {
	for _, part := range []string{"dataset.equivalentTo", "toggleEquivalentDiscoveredSubscription", "b.textContent=matched?'Unsubscribe':'Subscribe'", "recordAccountDiagnostic('Load'", "account-diagnostics-panel", "Browser subscriptions preserved"} {
		if !strings.Contains(personalJS,part) { t.Errorf("missing safeguard %q", part) }
	}
}
