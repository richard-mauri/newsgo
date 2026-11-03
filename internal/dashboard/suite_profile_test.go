package dashboard

import (
 "encoding/json"
 "testing"
)

func TestSuiteProfileLegacyRoundTrip(t *testing.T) {
 legacy:=storedPrefs{
  Revision:7, UpdatedAt:"2026-10-01T00:00:00Z",
  Subscriptions:[]syncSubscription{{ID:"feed:https://example.org/rss",Kind:"feed",Name:"Example",URL:"https://example.org/rss"}},
  ExcludedDomains:[]string{"example.com"},
  ExcludedArticleDomains:[]string{"ads.example.com"},
  MaxArticleAgeDays:7,MaximumLeadStories:8,LeadOnlyAllSubscriptions:true,
  LastReadingSelection:"https://example.org/rss",
 }
 data,err:=suiteDataFromLegacy(legacy);if err!=nil{t.Fatal(err)}
 encoded,err:=json.Marshal(data);if err!=nil{t.Fatal(err)}
 var shape map[string]any
 if err=json.Unmarshal(encoded,&shape);err!=nil{t.Fatal(err)}
 if _,ok:=shape["revision"];ok{t.Fatal("suite data must not contain legacy revision")}
 if _,ok:=shape["updatedAt"];ok{t.Fatal("suite data must not contain legacy timestamp")}
 restored,err:=legacyPrefsFromSuite(encoded,12);if err!=nil{t.Fatal(err)}
 if restored.Revision!=12||restored.MaximumLeadStories!=8||!restored.LeadOnlyAllSubscriptions||
 restored.LastReadingSelection!=legacy.LastReadingSelection||
 len(restored.Subscriptions)!=1||restored.Subscriptions[0].ID!=legacy.Subscriptions[0].ID||
 len(restored.ExcludedArticleDomains)!=1||restored.ExcludedArticleDomains[0]!="ads.example.com" {
  t.Fatalf("round trip changed settings: %+v",restored)
 }
}

func TestSuiteProfileRejectsMalformedData(t *testing.T) {
 if _,err:=legacyPrefsFromSuite(json.RawMessage(`{"subscriptions":[{"id":"bad","kind":"feed","name":"Bad","url":"file:///etc/passwd"}]}`),1);err==nil {t.Fatal("invalid subscription accepted")}
 if _,err:=legacyPrefsFromSuite(json.RawMessage(`{"maximumLeadStories":90}`),1);err==nil {t.Fatal("invalid lead count accepted")}
}
