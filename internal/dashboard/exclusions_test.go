package dashboard

import (
	"strings"
	"testing"
)

func TestExclusionControlsPreserveSubscriptionFlows(t *testing.T) {
	for _, required := range []string{
		"const key='newsgo.subscriptions.v2'",
		"function loadSubscriptions()",
		"function subscribeRecord(record)",
		"function unsubscribeRecord(id)",
		"function renderManager()",
		"function chooseAllSubscriptions()",
		"function renderPersonalArticles()",
		"sources.has(a.source)&&!excludedURL(a.url)",
		"function saveExclusions(list)",
		"function renderExclusionManager()",
		"function excludeButton(domain)",
		"function refreshExclusionViews()",
		"function showView(view)",
		"showView('news');searchCatalog();refresh();",
	} {
		if !strings.Contains(personalJS, required) {
			t.Errorf("missing subscription/exclusion contract: %s", required)
		}
	}
}
