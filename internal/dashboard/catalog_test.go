package dashboard

import (
	"strings"
	"testing"
)

func TestRemoteSearchAndSubscriptionPreservation(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="catalog-search"`, `id="catalog-list"`, `id="view-subscriptions"`, `id="catalog-search"`} {
		if !strings.Contains(h, s) {
			t.Fatalf("missing HTML %s", s)
		}
	}
	for _, s := range []string{"/api/search-feeds?q=", "searchCatalog", "subscribeClick({name:f.name", "localStorage.getItem(unifiedKey)", "submitCatalogSearch"} {
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
	for _, s := range []string{`showView('news')`, `showView('subscriptions')`, `syncFollowButtons()`, `textContent='Checking…'`, `following?(catalog?'Unsubscribe':'Subscribed'):b.dataset.subscribeError?'Failure':state==='unavailable'?'Unavailable':state==='unverified'?'Try Subscribe':'Subscribe'`, `textContent='('+total+')'`} {
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

func TestCatalogHidesBrokenFeedsAndTogglesSubscriptions(t *testing.T) {
 required := []string{
  "if(state==='unavailable'&&!followed){hiddenCount++;return}",
  "if(followed)subscribedCount++",
  "catalogSubscriptionClick(f,b)",
  "button(followed?'Unsubscribe'",
  "const catalog=b.closest('#catalog-list')!==null;",
  "following?(catalog?'Unsubscribe':'Subscribed')",
  "if(!await confirmSubscriptionRemoval(record.name))return;",
  "if(!unsubscribeRecord(record.id))return;",
  "localStorage.removeItem(readingSelectionKey);scheduleAccountSave();",
  "renderManager();renderReading();",
  "unavailable feeds hidden",
  "No usable feeds found",
 }
 for _,fragment := range required {
  if !strings.Contains(personalJS,fragment) {t.Errorf("missing filtered search/unsubscribe contract: %q",fragment)}
 }
}
