package dashboard

import (
 "encoding/json"
 "net/http"
 "strings"
 "testing"
)

func TestLastReadingSelectionServerSync(t *testing.T) {
 s, err := NewSyncStore(t.TempDir()); if err != nil {t.Fatal(err)}
 token, err := s.IssueToken("reader",""); if err != nil {t.Fatal(err)}
 feed := "https://example.com/feed.xml"
 prefs := storedPrefs{Subscriptions: []syncSubscription{{ID:"feed:"+feed,Kind:"feed",Name:"Example",URL:feed}}, LastReadingSelection:feed}
 w:=syncCall(t,s,http.MethodPut,token,prefs); if w.Code!=200 {t.Fatalf("save: %d %s",w.Code,w.Body.String())}
 w=syncCall(t,s,http.MethodGet,token,nil)
 var got userPrefs; if err:=json.Unmarshal(w.Body.Bytes(),&got);err!=nil {t.Fatal(err)}
 if got.LastReadingSelection!=feed || got.Revision!=1 {t.Fatalf("saved selection missing: %#v",got)}
 got.LastReadingSelection="all"
 w=syncCall(t,s,http.MethodPut,token,got.storedPrefs);if w.Code!=200 {t.Fatalf("save all: %d %s",w.Code,w.Body.String())}
 w=syncCall(t,s,http.MethodGet,token,nil);if err:=json.Unmarshal(w.Body.Bytes(),&got);err!=nil {t.Fatal(err)}
 if got.LastReadingSelection!="all" || got.Revision!=2 {t.Fatalf("all not saved: %#v",got)}
 got.Revision=1;got.LastReadingSelection=feed
 w=syncCall(t,s,http.MethodPut,token,got.storedPrefs);if w.Code!=409 {t.Fatalf("expected conflict, got %d",w.Code)}
}
func TestReadingSelectionValidationAndLegacyProfiles(t *testing.T){
 feed:="https://example.com/rss"
 p:=storedPrefs{Subscriptions:[]syncSubscription{{ID:"feed:"+feed,Kind:"feed",URL:feed,Name:"Example"}}}
 if !validPrefs(p){t.Fatal("legacy profile rejected")}
 for _,selection:=range []string{"all",feed}{p.LastReadingSelection=selection;if !validPrefs(p){t.Fatalf("valid selection rejected: %s",selection)}}
 for _,selection:=range []string{"https://unknown.example/rss","personal:"+feed,"javascript:alert(1)","ALL"}{p.LastReadingSelection=selection;if validPrefs(p){t.Fatalf("invalid selection accepted: %s",selection)}}
}
func TestReadingSelectionBrowserSyncContract(t *testing.T){
 for _,needle:=range []string{
  "lastReadingSelection:validReadingSelection(activeReadingSelection)||''",
  "restoreReadingSelection(validReadingSelection(data.lastReadingSelection))",
  "restoreReadingSelection(restored)",
  "rememberReadingSelection('all')",
  "rememberReadingSelection(url)",
  "snap.lastReadingSelection!==(validReadingSelection(activeReadingSelection)||'')",
  "else restoreReadingSelection(localStorage.getItem(readingSelectionKey))",
 } {if !strings.Contains(personalJS,needle){t.Errorf("missing browser restore contract %q",needle)}}
}
