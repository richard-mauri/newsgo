package dashboard

import ("strings";"testing")

func TestSingleCurrentSubscriptionRegistry(t *testing.T) {
 for _,token:=range []string{
  "const unifiedKey='newsgo.subscriptions.unified.v1'",
  "function loadSubscriptions()",
  "localStorage.getItem(unifiedKey)",
  "function saveSubscriptions(next)",
  "function subscribeRecord(record)",
  "function unsubscribeRecord(id)",
  "feedIdentity(f.url)",
  "snapshot.slice(offset,offset+8)",
 }{
  if !strings.Contains(personalJS,token){t.Errorf("missing registry feature %q",token)}
 }
 for _,obsolete:=range []string{
  "legacyPersonalKey","legacyIncludedKey","convertLegacy(","r.kind==='shared'",
  "newsgo.personalFeeds.v1","newsgo.includedSubscriptions.v1",
 }{if strings.Contains(personalJS,obsolete){t.Errorf("obsolete registry code remains: %q",obsolete)}}
}
