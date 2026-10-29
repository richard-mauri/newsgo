package dashboard

import (
 "strings"
 "testing"
)

func TestReadAllArticlesIsOnlyCombinedFeedNavigation(t *testing.T) {
 for _, fragment := range []string{
  "function chooseAllSubscriptions(leadOnly=false)", "el('read-source').value=''", "renderReading();showView('news')",
  "button('Read All Articles',()=>chooseAllSubscriptions(false))", "button('Read Lead Articles',()=>chooseAllSubscriptions(true))", "leadOnlyAllSubscriptions=leadOnly", "if(source&&excludedURL(publisherForFeed(source))){chooseAllSubscriptions();return}",
 } {
  if !strings.Contains(personalJS, fragment) { t.Errorf("missing combined-feed behavior: %s", fragment) }
 }
 for _, fragment := range []string{"id=\"all-subscriptions\"", "← All Subscriptions"} {
  if strings.Contains(page.Tree.Root.String(), fragment) { t.Errorf("obsolete standalone navigation control still present: %s",fragment) }
 }
 if strings.Contains(personalJS,"el('all-subscriptions')") { t.Error("obsolete all-subscriptions element still referenced in JS") }
 if !strings.Contains(page.Tree.Root.String(),`id="reading-title"`) {t.Error("selected publication title removed")}
}
