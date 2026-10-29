package dashboard

import (
	"strings"
	"testing"
)

func TestPersonalFeedsHaveLeadAndLatestSections(t *testing.T) {
	pageHTML := page.Tree.Root.String()
	for _, part := range []string{`id="personal-lead-section"`, `id="personal-lead-stories"`, `id="personal-latest-section"`, `id="personal-articles"`} {
		if !strings.Contains(pageHTML, part) {
			t.Errorf("missing personal headline section: %s", part)
		}
	}
	for _, part := range []string{`const featured=count++<maximumLeadStories;`, `if(!featured&&leadOnlyAllSubscriptions&&activeReadingSelection==='all')return;`, `el('personal-latest-section').hidden=!hasArticles||(leadOnlyAllSubscriptions&&activeReadingSelection==='all');`, `row.className=featured?'feature':'item'`, `lead.append(row)`, `list.append(row)`, `seen.has(a.url)`, `el('personal-lead-section').hidden=!hasArticles;`} {
		if !strings.Contains(personalJS, part) {
			t.Errorf("missing personal lead rendering: %s", part)
		}
	}
}
