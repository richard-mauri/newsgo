package dashboard

import (
	"strings"
	"testing"
)

// Ensure the all-subscriptions branch includes selected shared publishers and
// does not revert to the v0.2.23 always-hidden shared headline behavior.
func TestCombinedSubscriptionsVisibility(t *testing.T) {
	for _, expected := range []string{
		"included.includes(a.dataset.source)",
		"source===''&&included.length>0",
		"source===''&&feeds.length>0",
		"el('subscriptions-empty-notice').hidden=source!==''||registry.length>0",
	} {
		if !strings.Contains(personalJS, expected) {
			t.Fatalf("missing combined-view behavior: %s", expected)
		}
	}
}
