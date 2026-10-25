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
  "window.confirm(",
  "unsubscribeRecord(record.id)",
  "chooseAllSubscriptions()",
  "renderFilteredSubscriptionOffer();",
  "panel.hidden=true",
  "title.insertAdjacentElement('afterend',panel)",
 }
 for _,fragment:=range required {
  if !strings.Contains(personalJS,fragment) {t.Errorf("missing UI contract: %q",fragment)}
 }
}
