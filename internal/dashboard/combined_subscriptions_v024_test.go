package dashboard

import (
	"strings"
	"testing"
)

func TestUnifiedSubscriptionVisibility(t *testing.T) {
	for _, s := range []string{"function readingSelection(source)", "showPersonal:true", "function renderPersonalArticles()", "sources.has(a.source)&&!excludedURL(a.url)", "function convertLegacy(records)", "catalogByName.get(r.name)"} {
		if !strings.Contains(personalJS, s) {
			t.Errorf("missing unified subscription behavior: %s", s)
		}
	}
	for _, s := range []string{"function arrangeSharedLeads(selection)", "showShared:", "included.includes(name)"} {
		if strings.Contains(personalJS, s) {
			t.Errorf("obsolete independent shared headline pipeline: %s", s)
		}
	}
}
