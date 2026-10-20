package dashboard

import (
	"strings"
	"testing"
)

func TestRemoteSearchAndSubscriptionPreservation(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="catalog-search"`, `id="catalog-list"`, `id="view-subscriptions"`, `/personal.js?v=0.2.41`} {
		if !strings.Contains(h, s) {
			t.Fatalf("missing HTML %s", s)
		}
	}
	for _, s := range []string{"/api/search-feeds?q=", "searchCatalog", "add({name:f.name", "localStorage.getItem(key)", "queueSearch"} {
		if !strings.Contains(personalJS, s) {
			t.Fatalf("missing JS %s", s)
		}
	}
	if strings.Contains(personalJS, "const catalog=[") {
		t.Fatal("static-only catalog remains")
	}
}
func TestUXNavigationAndSubscriptionManagement(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="manage-personal"`, `id="view-subscriptions"`, `id="subscriptions-panel"`, `id="personal-sources"`, `id="close-discovery"`, `id="close-subscriptions"`, `id="news-content"`} {
		if !strings.Contains(h, s) {
			t.Fatalf("missing navigation element %s", s)
		}
	}
	for _, s := range []string{`showView('news')`, `showView('subscriptions')`, `syncFollowButtons()`, `textContent='Checking…'`, `following?'Subscribed':'Subscribe'`, `textContent='('+total+')'`} {
		if !strings.Contains(personalJS, s) {
			t.Fatalf("missing subscription interaction %s", s)
		}
	}
}
func TestPersonalArticleSourceSelection(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="read-source"`, `id="followed-options"`, `id="reading-title"`, `Select Read Articles to view headlines`} {
		if !strings.Contains(h, s) {
			t.Fatalf("missing personal source UI %s", s)
		}
	}
	for _, s := range []string{`function chooseSource(url)`, `function renderPersonalArticles()`, `renderReading()`, `sources.has(a.source)`, `showView('news')`} {
		if !strings.Contains(personalJS, s) {
			t.Fatalf("missing personal source selection logic %s", s)
		}
	}
}
