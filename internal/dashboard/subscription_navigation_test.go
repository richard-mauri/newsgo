package dashboard

import (
	"strings"
	"testing"
)

func TestSubscriptionFirstNavigation(t *testing.T) {
	html := page.Tree.Root.String()
	for _, item := range []string{`Find Subscriptions`, `My Subscriptions`, `Read Articles`, `Unsubscribe`, `id="reading-title"`} {
		if !strings.Contains(html, item) && !strings.Contains(personalJS, item) {
			t.Errorf("missing subscription UX element: %s", item)
		}
	}
	for _, old := range []string{`id="view-news"`, `id="clear-reading"`, `id="read-category"`, `id="back-subscriptions"`, `Back to My Subscriptions`, `Unregister`} {
		if strings.Contains(html, old) {
			t.Errorf("obsolete control: %s", old)
		}
	}
	for _, js := range []string{`newsgo.subscriptions.v2`, `newsgo.personalFeeds.v1`, `newsgo.includedSubscriptions.v1`, `saveSubscriptions(`, `chooseIncluded(name)`, `button('Unsubscribe'`, `button('Read Articles'`} {
		if !strings.Contains(personalJS, js) {
			t.Errorf("missing subscription logic: %s", js)
		}
	}
}
