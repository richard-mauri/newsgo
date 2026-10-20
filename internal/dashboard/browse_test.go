package dashboard

import (
	"strings"
	"testing"
)

func TestBrowseIncludedAndPublisherOnly(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="read-source"`, `data-feed-url="{{.URL}}"`, `Browse available sources`, `Subscribe`, `/personal.js?v=0.2.32`} {
		if !strings.Contains(h, s) {
			t.Fatalf("missing browse UI: %s", s)
		}
	}
	if strings.Contains(h, `id="read-category"`) || strings.Contains(personalJS, "read-category") {
		t.Fatal("obsolete category filtering present")
	}
	for _, s := range []string{`.read-included`, `saveSubscriptions(`, `chooseIncluded(name)`} {
		if !strings.Contains(personalJS, s) {
			t.Fatalf("missing included source navigation: %s", s)
		}
	}
}
