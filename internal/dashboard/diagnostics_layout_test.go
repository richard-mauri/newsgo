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
	if strings.Contains(html[nav:diagnostics], `id="exclusions-panel"`) {
		t.Fatal("Excluded Domains disclosure must not be on the main page")
	}
	panel := html[diagnostics:diagnosticsEnd]
	for _, id := range []string{
		`<details id="exclusions-panel"`,
		`<summary><strong>Excluded Domains</strong> <span id="excluded-count">(0)</span></summary>`,
		`id="filtered-articles-card"`,
		`id="filtered-articles-list"`,
		`id="filtered-articles-count"`,
		`<details id="raw-feed-inspector" class="personal-panel" aria-label="Raw feed responses"><summary><strong>Raw feed responses</strong> <span id="raw-feed-count"></span></summary>`,
		`id="raw-feed-list"`,
		`id="raw-feed-status"`,
	} {
		if !strings.Contains(panel, id) {
			t.Errorf("Diagnostics missing %s", id)
		}
	}
	if strings.Contains(panel, `id="manage-exclusions"`) || strings.Contains(panel, `id="close-exclusions"`) {
		t.Fatal("obsolete Excluded Domains navigation controls remain")
	}
	if strings.Contains(panel, `<section id="raw-feed-inspector"`) || strings.Contains(panel, `<h3>Raw feed responses`) {
		t.Fatal("raw feed inspector still uses non-disclosure heading")
	}
	if strings.Contains(html[news:], `id="filtered-articles-card"`) {
		t.Fatal("filtered audit still appears in main article content")
	}
}

func TestDiagnosticsExclusionNavigationAndFilteringRetained(t *testing.T) {
	for _, marker := range []string{
		"el('diagnostics-panel').classList.toggle('personal-hidden',view!=='diagnostics')",
		"el('manage-diagnostics').addEventListener('click',()=>{showView('diagnostics');updateDiagnostics()})",
		"function renderFilteredSubscriptionOffer()",
		"function renderFilteredArticles()",
		"renderFilteredArticles();",
		"function setUnifiedDomain(domain,allow)",
		"const details=document.createElement('details');details.className='raw-feed-entry'",
		"const download=button('Download Raw Feed',()=>{",
		"const copy=button('Copy Raw Feed',async()=>{",
		"list.replaceChildren();count.textContent='('+latestRawFeeds.length+' feeds)'",
		"el('excluded-count').textContent='('+domains.length+')",
	} {
		if !strings.Contains(personalJS, marker) {
			t.Errorf("missing diagnostics/exclusions behavior: %s", marker)
		}
	}
}
