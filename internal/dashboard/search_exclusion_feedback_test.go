package dashboard

import (
 "strings"
 "testing"
)

// Regression: a successful exclusion must update every visible search candidate
// for that publisher without requiring another search or a page reload.
func TestSearchCandidateExclusionFeedback(t *testing.T) {
 for _, marker := range []string{
  "renderExclusionManager();renderManager();refreshSearchExclusionCards();",
  "function refreshSearchExclusionCards(){",
  "document.querySelectorAll('#catalog-list .catalog-card')",
  "action.dataset.publisherDomain",
  "action.textContent=excluded?'Allow this site':'Exclude this site'",
  "card.dataset.excluded=excluded?'true':'false'",
  "badge.textContent=statusLabel(domain)",
  "b.dataset.publisherDomain=domain;b.disabled=followed||excludedURL(domain)",
  "b.disabled=following||excludedURL(domain)",
  "Allow this site before subscribing",
  "syncFollowButtons();",
  "function publisherScope(value){",
  "return d.endsWith('.gov.au')&&parts.length>3?parts.slice(-3).join('.'):d;",
  "if(excludedURL(domain))saveExclusions(excludedDomains.filter(x=>publisherScope(x)!==domain))",
  "function publisherForFeed(feed){return publisherScope(feed.website||feed.url||'')}",
 } {
  if !strings.Contains(personalJS, marker) { t.Errorf("missing search exclusion refresh contract: %s", marker) }
 }
}
