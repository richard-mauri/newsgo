package dashboard

import (
 "strings"
 "testing"
)

func TestFilteredSubscriptionOfferContract(t *testing.T) {
 required:=[]string{
  "activeReadingSelection==='all'",
  "evaluated>=10&&rejected/evaluated>0.9",
  "'too-old','publisher exclusion','article-domain exclusion'",
  "confirmSubscriptionRemoval(feed.name)",
  "sync-choice-backdrop",
  "sync-choice-dialog",
  "exclusion-confirm-action",
  "if(!rememberedSyncToken)renderPersonalArticles()",
  "unsubscribeRecord(record.id)",
  "chooseAllSubscriptions()",
  "renderFilteredSubscriptionOffer();",
  "panel.hidden=true",
  "heading.insertAdjacentElement('afterend',panel)",
  "anchor.insertAdjacentElement('beforebegin',heading)",
  "heading.textContent=feed.name",
  "el('reading-title').hidden=true",
  "action.className='exclusion-confirm-action'",
 }
 for _,fragment:=range required {
  if !strings.Contains(personalJS,fragment) {t.Errorf("missing UI contract: %q",fragment)}
 }
}
