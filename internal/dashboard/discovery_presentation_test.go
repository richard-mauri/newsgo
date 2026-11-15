package dashboard

import (
	"net/url"
	"strings"
	"testing"
)

func TestReadableDiscoveryPresentation(t *testing.T) {
	u, err := url.Parse("https://www.foxsports.com/college-football")
	if err != nil {
		t.Fatal(err)
	}
	got := discoveryName([]byte(`<html><title>College Football | FOX Sports</title></html>`), u, "https://api.foxsports.com/v2/content/optimized-rss?tags=fs%2Fcfb")
	if !strings.Contains(got, "College Football") || strings.Contains(got, "partnerKey") {
		t.Fatalf("unfriendly title: %q", got)
	}
	if c := discoveryCategory(got, u.String()); c != "Sports" {
		t.Fatalf("category %q", c)
	}
	if !strings.Contains(personalJS, "discovered-feed-address") || !strings.Contains(personalJS, "address.textContent=f.url") || !strings.Contains(personalJS, "data.counts") {
		t.Fatal("missing friendly result fields")
	}
	if strings.Contains(personalJS, "summary.textContent='Technical feed address'") {
		t.Fatal("discovered feed URLs must not be collapsed")
	}
	for _, part := range []string{"discovered-feed-row", "subscribeButton.dataset.feedUrl=f.url", "toggleDiscoveredSubscription(f,subscribeButton)", "toggleDiscoveredSubscription(f,b)", "if(feeds.some(existing=>existing.url===f.url)){catalogSubscriptionClick(f,trigger);return}", "if(!subscribeRecord(feedRecord(f)))throw Error"} {
		if !strings.Contains(personalJS, part) { t.Fatalf("missing discovery subscribe state handling: %s", part) }
	}
	if !strings.Contains(page.Tree.Root.String(), "grid-template-columns:minmax(0,1fr) 128px 205px") {
		t.Fatal("discovery subscribe and exclude controls must have aligned columns")
	}
}

func TestDiscoveryFeedMetadataAndDistinctFallbacks(t *testing.T) {
	u, _ := url.Parse("https://www.investing.com/webmaster-tools/rss")
	page := []byte(`<html><title>Investing.com RSS Feeds</title></html>`)
	stockURL := "https://www.investing.com/rss/stock.rss"
	forexURL := "https://www.investing.com/rss/forex.rss"
	stock := discoveryFeedName(page, []byte(`<rss><channel><title>Investing.com Stock Market News</title></channel></rss>`), u, stockURL)
	if stock != "Investing.com Stock Market News" { t.Fatalf("RSS feed title: %q", stock) }
	forex := discoveryFeedName(page, []byte(`<rss><channel><title>RSS Feeds</title></channel></rss>`), u, forexURL)
	if forex == stock || !strings.Contains(forex, "Forex") { t.Fatalf("missing distinct fallback: %q", forex) }
	if c := discoveryCategory(forex, u.String()+" "+forexURL); c != "Business" { t.Fatalf("forex category: %q", c) }
	if c := discoveryCategory(stock, u.String()+" "+stockURL); c != "Business" { t.Fatalf("stocks category: %q", c) }
	atom := discoveryFeedName(page, []byte(`<feed xmlns="http://www.w3.org/2005/Atom"><title>World Science</title></feed>`), u, "https://www.investing.com/atom.xml")
	if atom != "World Science" { t.Fatalf("Atom title: %q", atom) }
	json := discoveryFeedName(page, []byte(`{"version":"https://jsonfeed.org/version/1.1","title":"Technology Briefing"}`), u, "https://www.investing.com/feed.json")
	if json != "Technology Briefing" { t.Fatalf("JSON title: %q", json) }
}
func TestVersionedDiscoveryScript(t *testing.T) {
	assertVersionedScript(t,page.Tree.Root.String())
}

func TestAlignedCatalogActionsAndSharedButtonFeedback(t *testing.T) {
	style := page.Tree.Root.String()
	for _, want := range []string{
		".included-card{display:grid;grid-template-columns:minmax(0,1fr) 128px 205px",
		".included-card>.read-included{grid-column:2;grid-row:1}",
		".included-card>.exclude-site{grid-column:3;grid-row:1}",
		".included-card>.exclusion-label{grid-column:1;grid-row:2",
		".catalog-card{display:grid;grid-template-columns:minmax(0,1fr) 128px 205px",
		".included-card strong{grid-column:1/-1;grid-row:1}",
		".personal-row.discovered-feed-row{display:grid;grid-template-columns:minmax(0,1fr) 128px 205px",
		".included-card button{width:100%;min-width:0;cursor:pointer;background:var(--soft)",
		"button:not(:disabled):hover{color:var(--accent);border-color:var(--accent)",
		"button:not(:disabled):focus-visible{outline:2px solid var(--accent)",
		"@media(prefers-reduced-motion:reduce){button:not(:disabled){transition:none}}",
	} {
		if !strings.Contains(style, want) {
			t.Errorf("missing aligned or accessible button treatment: %s", want)
		}
	}
}
