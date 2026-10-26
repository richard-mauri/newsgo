package dashboard

import (
 "strings"
 "testing"
)

func TestFilteredSubscriptionOfferContract(t *testing.T) {
 required:=[]string{
  "activeReadingSelection==='all'",
  "total>0&&filtered/total>0.9",
  "const classifiedFiltered=filteredEntries.filter(a=>a.source===sourceName).length",
  "const excluded=articles.filter(a=>excludedURL(a.url)||articleExcluded(a.url)).length",
  "confirmSubscriptionRemoval(feed.name)",
  "sync-choice-backdrop",
  "sync-choice-dialog",
  "exclusion-confirm-action",
  "if(!rememberedSyncToken)renderPersonalArticles()",
  "unsubscribeRecord(record.id)",
  "activeReadingSelection=null;selectedSource=''",
  "Subscription removed. Choose another subscription to read.",
  "renderFilteredSubscriptionOffer();",
  "panel.hidden=true",
  "heading.insertAdjacentElement('afterend',panel)",
  "anchor.insertAdjacentElement('beforebegin',heading)",
  "heading.textContent=feed.name",
  "el('reading-title').hidden=true",
  "action.className='subscription-unsubscribe-button'",
  "select.value=activeReadingSelection",
  "const restored=activeReadingSelection",
 }
 for _,fragment:=range required {
  if !strings.Contains(personalJS,fragment) {t.Errorf("missing UI contract: %q",fragment)}
 }
}
