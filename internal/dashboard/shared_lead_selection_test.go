package dashboard

import (
	"strings"
	"testing"
)

func TestUnifiedLeadSelectionHonorsExclusions(t *testing.T) {
	for _, s := range []string{"const articles=allPersonalArticles.filter(a=>sources.has(a.source)&&!excludedURL(a.url));", "const featured=count++<3;", "lead.append(row)", "list.append(row)"} {
		if !strings.Contains(personalJS, s) {
			t.Errorf("missing lead selection: %s", s)
		}
	}
}
