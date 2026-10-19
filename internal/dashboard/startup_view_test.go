package dashboard

import (
	"strings"
	"testing"
)

func TestSubscriptionFirstStartupDoesNotFlashSharedHeadlines(t *testing.T) {
	if !strings.Contains(page.Tree.Root.String(), `id="shared-headlines" class="personal-hidden"`) {
		t.Fatal("shared headline cards must be hidden in initial HTML to prevent first-paint flash")
	}
	if !strings.Contains(personalJS, "const showShared=source.startsWith('shared:');") {
		t.Fatal("shared headline cards should appear only when a built-in publisher is explicitly selected")
	}
	if !strings.Contains(page.Tree.Root.String(), `/personal.js?v=0.2.23`) {
		t.Fatal("browser JS URL must be cache-busted for startup fix")
	}
}
