package dashboard

import (
	"strings"
	"testing"
)

func TestBrowseIncludedAndPublisherOnly(t *testing.T) {
	h := page.Tree.Root.String()
	for _, expected := range []string{`id="read-source"`, `data-included-source="{{.}}"`, `Browse included sources`, `Subscribe`, `/personal.js?v=0.2.23`} {
		if !strings.Contains(h, expected) {
			t.Fatalf("missing browse UI: %s", expected)
		}
	}
	if strings.Contains(h, `id="read-category"`) {
		t.Fatal("obsolete section selector present")
	}
	if strings.Contains(personalJS, "read-category") {
		t.Fatal("obsolete category filtering present")
	}
	for _, expected := range []string{`.read-included`, `saveSubscriptions(`, `chooseIncluded(name)`} {
		if !strings.Contains(personalJS, expected) {
			t.Fatalf("missing included source navigation: %s", expected)
		}
	}
}
