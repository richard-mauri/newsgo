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
	if strings.Contains(personalJS, "accountSnapshot(){return {recentHistory") {
		t.Fatal("recent discovery history must not be included in Homepage account storage")
	}
}
