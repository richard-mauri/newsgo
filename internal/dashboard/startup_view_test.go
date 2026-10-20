package dashboard

import (
	"strings"
	"testing"
)

func TestSubscriptionFirstStartupHasNoUnsolicitedArticles(t *testing.T) {
	h := page.Tree.Root.String()
	for _, s := range []string{`id="personal-feed"`, `id="personal-articles"`, `/personal.js?v=0.2.32`} {
		if !strings.Contains(h, s) {
			t.Errorf("missing unified startup UI %s", s)
		}
	}
	for _, s := range []string{`id="shared-headlines"`, `{{range $i,$a := .Articles}}`, `cached headlines from trusted feeds`} {
		if strings.Contains(h, s) {
			t.Errorf("obsolete automatic headlines: %s", s)
		}
	}
}
