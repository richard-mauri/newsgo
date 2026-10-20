package dashboard

import (
	"strings"
	"testing"
)

func TestNoRedundantBackPanelAndConsistentTerms(t *testing.T) {
	html := page.Tree.Root.String()
	if strings.Contains(html, `id="back-subscriptions"`) {
		t.Fatal("redundant back button")
	}
	if !strings.Contains(html, `id="reading-title"`) {
		t.Fatal("selected publication heading missing")
	}
	for _, old := range []string{"Unregister", "Unfollow", "Back to My Subscriptions"} {
		if strings.Contains(html, old) || strings.Contains(personalJS, old) {
			t.Errorf("old terminology: %s", old)
		}
	}
	if !strings.Contains(personalJS, `button('Unsubscribe'`) {
		t.Fatal("unsubscribe action missing")
	}
}
