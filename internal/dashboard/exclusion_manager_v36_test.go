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
        "function setUnifiedDomain(domain,allow)",
        "const domains=effectiveExcludedDomains().sort()", 
        "row.append(info,button('Allow domain',()=>setUnifiedDomain(d,true)))", 
        "el('excluded-count').textContent='('+domains.length+')'", 
        "()=>{renderExclusionManager();showView('exclusions')}",
    } {
        if !strings.Contains(personalJS, token) { t.Errorf("missing exclusion manager contract %q", token) }
    }
}
