package dashboard

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync/atomic"
 "testing"
 "time"
)

func TestReaderCachePreferenceBypassesWithoutClearingSharedEntry(t *testing.T) {
 clearNewsClassificationCache()
 t.Cleanup(clearNewsClassificationCache)
 var calls atomic.Int32
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/classify" {t.Errorf("unexpected endpoint %s",r.URL.Path)}
  if calls.Add(1)==1 {_,_=w.Write([]byte(`{"classification":"article"}`));return}
  _,_=w.Write([]byte(`{"classification":"homepage"}`))
 }))
 defer server.Close()
 article:="https://example.org/newsgo-cache-preference-case"
 // Explicitly prime shared server cache to avoid dependence on admin defaults.
 normalized:=normalizeClassificationURL(article)
 key:=server.URL+"|"+normalized
 newsClassificationCache.Lock()
 newsClassificationCache.entries[key]=newsClassificationEntry{value:"article",expiry:time.Now().Add(time.Hour)}
 newsClassificationCache.Unlock()
 if got:=newsClassifyWithCache(context.Background(),server.URL,article,true);got!="article"{t.Fatalf("cache enabled got %q",got)}
 if got:=calls.Load();got!=0{t.Fatalf("cache hit contacted ScraperGo %d times",got)}
 if got:=newsClassifyWithCache(context.Background(),server.URL,article,false);got!="article"{t.Fatalf("fresh result got %q",got)}
 if calls.Load()!=1{t.Fatalf("bypass did not contact ScraperGo")}
 // Bypass must not replace the shared cached value used by other readers.
 if got:=newsClassifyWithCache(context.Background(),server.URL,article,true);got!="article"{t.Fatalf("bypass overwrote cache: %q",got)}
 if calls.Load()!=1{t.Fatalf("unexpected extra classification request")}
}

func TestDefaultOnCacheRequestAndLegacyProfile(t *testing.T) {
 var req personalRequest
 if err:=json.Unmarshal([]byte(`{"feeds":[]}`),&req);err!=nil{t.Fatal(err)}
 if req.UseClassificationCache!=nil{t.Fatal("old request must default to use-cache")}
 if err:=json.Unmarshal([]byte(`{"feeds":[],"useClassificationCache":false}`),&req);err!=nil{t.Fatal(err)}
 if req.UseClassificationCache==nil||*req.UseClassificationCache{t.Fatal("explicit bypass lost")}

}

func TestClassificationCachePreferenceWiring(t *testing.T){
 for _,needle:=range []string{
  `id="use-classification-cache" checked`,
  "Use classification cache",
 }{
  if !strings.Contains(page.Tree.Root.String(),needle){t.Errorf("UI missing %q",needle)}
 }
 for _,needle:=range []string{
  "localStorage.getItem(classificationCachePreferenceKey)!=='false'",
  "JSON.stringify({feeds:batch,maxArticleAgeDays,rssMetadataOnly,useClassificationCache})",
  "JSON.stringify({feeds:[f],maxArticleAgeDays,rssMetadataOnly,useClassificationCache})",
  "useClassificationCache",
  "useClassificationCache,lastReadingSelection:",
  "data.useClassificationCache!==false",
  "scheduleServerSave();",
 }{
  if !strings.Contains(personalJS,needle){t.Errorf("browser wiring missing %q",needle)}
 }
}
