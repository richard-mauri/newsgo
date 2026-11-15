package dashboard

import (
	"strings"
	"testing"
)

func TestRecentDiscoveryHistoriesRemainSeparateAndBrowserLocal(t *testing.T) {
	for _, part := range []string{
		"newsgo.recent.searches.v1",
		"newsgo.recent.websites.v1",
		"recentHistoryLimit=8",
		"rememberRecentHistory(recentHistoryFields[0],q)",
		"rememberRecentHistory(recentHistoryFields[1],u)",
		"field.run()",
		"Clear history",
		"localStorage.removeItem(field.key)",
	} {
		if !strings.Contains(personalJS, part) {
			t.Errorf("missing recent-history behavior %q", part)
		}
	}
	if !strings.Contains(personalJS,"discoveryFeedback('Could not find feeds: '+err.message)") || !strings.Contains(personalJS,"Enter a valid public http:// or https:// website URL.") {
		t.Fatal("invalid website discovery must show an inline error")
	}
	validResult := strings.Index(personalJS,"if(!data.feeds||!data.feeds.length)throw Error")
	remember := strings.Index(personalJS,"rememberRecentHistory(recentHistoryFields[1],u)")
	if validResult<0 || remember<=validResult {
		t.Fatal("website URLs must enter history only after valid feed discovery")
	}
}
