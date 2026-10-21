package dashboard

import (
	"strings"
	"testing"
)

func TestManualSearchSubmission(t *testing.T) {
	html := page.Tree.Root.String()
	for _, s := range []string{`id="catalog-form"`, `id="catalog-submit"`, `type="submit"`} {
		if !strings.Contains(html, s) {
			t.Errorf("missing search HTML %q", s)
		}
	}
	for _, s := range []string{
		"addEventListener('submit',submitCatalogSearch)",
		"addEventListener('input',queryEdited)",
		"Press Search or Enter", "res.status!==429", "Previous results remain visible", "searchIntervalMs=3200",
	} {
		if !strings.Contains(personalJS, s) {
			t.Errorf("missing manual search behavior %q", s)
		}
	}
	for _, bad := range []string{
		"setTimeout(searchCatalog,350)", "addEventListener('input',queueSearch)",
		"showView('news');searchCatalog();refresh()",
	} {
		if strings.Contains(personalJS, bad) {
			t.Errorf("unexpected automatic search: %q", bad)
		}
	}
}
