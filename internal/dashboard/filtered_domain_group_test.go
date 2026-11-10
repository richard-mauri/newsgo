package dashboard

import (
 "strings"
 "testing"
)

// The filtered audit must not render one exclusion control per article.
// Grouping by the actual URL domain also avoids trusting feed titles as hostnames.
func TestFilteredAuditGroupsPublisherDomains(t *testing.T) {
 script := personalJS
 for _, want := range []string{
  "const byDomain=new Map()",
  "const domain=publisherScope(article.url)",
  "for(const [domain,domainItems] of byDomain)",
  "name.textContent=(domain||'Unknown domain')+' ('+domainItems.length+')'",
  "if(domain&&!domainIsExcluded(domain))",
  "const exclude=button('Exclude '+domain,async()=>{",
  "await setUnifiedDomain(domain,false)",
  "group.append(row)",
  "section.append(group)",
  "summary.textContent='Admin diagnostics (run fresh classification)'",
 } {
  if !strings.Contains(script,want) {t.Errorf("missing grouped audit invariant %q",want)}
 }
 if strings.Index(script,"const byDomain=new Map()") > strings.Index(script,"for(const a of domainItems)") {
  t.Fatal("domain group must be constructed before individual rows")
 }
 for _, wrong := range []string{
  "excludeDomain(a.url)",
  "setUnifiedDomain(a.url,false)",
 } {
  if strings.Contains(script,wrong) {t.Errorf("per-article exclusion instead of per-domain: %q",wrong)}
 }
}

func TestFilteredDomainControlAndExistingConfirmation(t *testing.T) {
 markup:=page.Tree.Root.String()
 for _,want:=range []string{
  ".filtered-domain-group{",
  ".filtered-domain-heading{",
  ".filtered-domain-exclude{",
 } {
  if !strings.Contains(markup,want) {t.Errorf("missing grouped audit CSS %q",want)}
 }
 for _,want:=range []string{
  "function publisherScope(value)",
  "function confirmDomainRemoval(domain,affected)",
  "if(!allow&&!(await confirmDomainRemoval(domain,affected)))return false",
  "scheduleAccountSave();refresh()",
  "if(domain&&!domainIsExcluded(domain))",
 } {
  if !strings.Contains(personalJS,want) {t.Errorf("lost safe exclusion behavior %q",want)}
 }
}
