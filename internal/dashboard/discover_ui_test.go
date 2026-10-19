package dashboard

import (
	"strings"
	"testing"
)

func TestDiscoveryButtonDoesNotSubmitPage(t *testing.T) {
	if !strings.Contains(page.Tree.Root.String(), `id="discover-button" type="button"`) {
		t.Fatal("Find feeds must not submit the page")
	}
	if !strings.Contains(personalJS, "addEventListener('click',discoverFeeds)") {
		t.Fatal("missing click handler")
	}
	if !strings.Contains(personalJS, "e.preventDefault();discoverFeeds()") {
		t.Fatal("missing Enter-key form handler")
	}
}
