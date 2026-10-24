package dashboard

import (
    "strings"
    "testing"
)

func TestExcludedManagerNamesSitesAndPreservesDomainCountV36(t *testing.T) {
    // A saved exclusion such as news.ycombinator.com must appear by its
    // subscription name, even though the persisted key is a domain.
    for _, token := range []string{
        "function excludedURLForDomain(url,domain)",
        "registry.filter(r=>r.kind==='feed'&&publisherScope(r.url)===publisherScope(d))",
        "names.join(' · ')",
        "host.textContent=d",
        "el('excluded-count').textContent='('+(excludedDomains.length+articleDomains.length)+')'",
        "()=>{renderExclusionManager();showView('exclusions')}",
    } {
        if !strings.Contains(personalJS, token) { t.Errorf("missing exclusion manager contract %q", token) }
    }
}
