package dashboard

import (
	"strings"
	"testing"
)

func TestSubscriptionFirstStartupDoesNotFlashSharedHeadlines(t *testing.T) {
	if !strings.Contains(page.Tree.Root.String(), `id="shared-headlines" class="personal-hidden"`) {
		t.Fatal("shared headlines must be hidden until subscriptions load")
	}
	for _, s := range []string{"showShared:!!sharedName||(all&&included.length>0)", "included.includes(name)"} {
		if !strings.Contains(personalJS, s) {
			t.Fatalf("shared headlines must appear only for subscribed publishers: %s", s)
		}
	}
	if !strings.Contains(page.Tree.Root.String(), `/personal.js?v=0.2.31`) {
		t.Fatal("browser JS URL must be versioned")
	}
}
