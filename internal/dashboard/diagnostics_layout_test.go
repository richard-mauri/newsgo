package dashboard

import (
	"strings"
	"testing"
)

func TestDiagnosticsOwnsFilteredAuditAndExcludedDomains(t *testing.T) {
	html := page.Tree.Root.String()
	nav := strings.Index(html, `<nav class="news-nav"`)
	diagnostics := strings.Index(html, `id="diagnostics-panel"`)
	diagnosticsEnd := strings.Index(html, `id="homepage-confirm-dialog"`)
	news := strings.Index(html, `id="news-content"`)
	if nav < 0 || diagnostics < 0 || diagnosticsEnd <= diagnostics || news <= diagnosticsEnd {
		t.Fatal("unexpected dashboard section ordering")
	}
	if strings.Contains(html[nav:diagnostics], `id="manage-exclusions"`) {
		t.Fatal("Excluded Domains navigation button must not be on the main page")
	}
	panel := html[diagnostics:diagnosticsEnd]
	for _, id := range []string{
		`id="manage-exclusions"`,
		`id="exclusions-panel"`,
		`id="filtered-articles-card"`,
		`id="filtered-articles-list"`,
		`id="filtered-articles-count"`,
	} {
		if !strings.Contains(panel, id) {
			t.Errorf("Diagnostics missing %s", id)
		}
	}
	if strings.Contains(html[news:], `id="filtered-articles-card"`) {
		t.Fatal("filtered audit still appears in main article content")
	}
}

func TestDiagnosticsExclusionNavigationAndFilteringRetained(t *testing.T) {
	for _, marker := range []string{
		"view!=='diagnostics'&&view!=='exclusions'",
		"el('manage-exclusions').addEventListener('click',()=>{renderExclusionManager();showView('exclusions')})",
		"el('close-exclusions').addEventListener('click',()=>showView('diagnostics'))",
		"el('manage-diagnostics').addEventListener('click',()=>{showView('diagnostics');updateDiagnostics()})",
		"function renderFilteredSubscriptionOffer()",
		"function renderFilteredArticles()",
		"renderFilteredArticles();",
		"function setUnifiedDomain(domain,allow)",
		"el('excluded-count').textContent='('+domains.length+')",
	} {
		if !strings.Contains(personalJS, marker) {
			t.Errorf("missing diagnostics/exclusions behavior: %s", marker)
		}
	}
}
