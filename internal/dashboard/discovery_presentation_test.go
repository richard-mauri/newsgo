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
	if !strings.Contains(personalJS, "Technical feed address") || !strings.Contains(personalJS, "data.counts") {
		t.Fatal("missing friendly result fields")
	}
}
func TestVersionedDiscoveryScript(t *testing.T) {
	assertVersionedScript(t,page.Tree.Root.String())
}
