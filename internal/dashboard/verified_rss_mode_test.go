package dashboard

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "sync/atomic"
 "testing"
 "time"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

// RSS-only changes presentation and image choice, never the article filter.
func TestRSSMetadataOnlyStillFiltersNonArticlesAndUsesClassificationCache(t *testing.T) {
 previous:=currentArticleSettings()
 defer activeArticleSettings.Store(previous)
 activeArticleSettings.Store(defaultArticleSettings())
 clearNewsClassificationCache()
 defer clearNewsClassificationCache()
 source:=feeds.Source{Name:"Verified RSS",Category:"Other",URL:"https://example.org/verified-rss.xml"}
 now:=time.Now()
 personalCache.Lock()
 previousFeed,exists:=personalCache.entries[source.URL]
 personalCache.entries[source.URL]=cacheEntry{articles:[]feeds.Article{
  {Title:"Index landing page",URL:"https://example.org/index-test",Published:now},
  {Title:"Real story",URL:"https://example.org/article-test",Published:now,Image:"https://example.org/thumbnail.jpg"},
 },at:now}
 personalCache.Unlock()
 defer func(){personalCache.Lock();if exists{personalCache.entries[source.URL]=previousFeed}else{delete(personalCache.entries,source.URL)};personalCache.Unlock()}()
 var calls atomic.Int32
 reader:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/classify"{t.Errorf("unexpected path %s",r.URL.Path);return}
  calls.Add(1)
  var v struct{URL string `json:"url"`}
  if err:=json.NewDecoder(r.Body).Decode(&v);err!=nil{t.Error(err)}
  w.Header().Set("Content-Type","application/json")
  if v.URL=="https://example.org/index-test"{_,_=w.Write([]byte(`{"classification":"homepage"}`));return}
  if v.URL=="https://example.org/article-test"{_,_=w.Write([]byte(`{"classification":"article"}`));return}
  t.Errorf("unexpected classification URL %q",v.URL)
 }))
 defer reader.Close()
 state:=&State{reader:reader.URL}
 get:=func(rss bool,enabled bool) personalResult{
  t.Helper()
  body,_:=json.Marshal(personalRequest{Feeds:[]feeds.Source{source},RSSMetadataOnly:rss,UseClassificationCache:&enabled})
  w:=httptest.NewRecorder()
  state.servePersonal(w,httptest.NewRequest(http.MethodPost,"/api/personal",bytes.NewReader(body)))
  if w.Code!=http.StatusOK{t.Fatalf("status %d %s",w.Code,w.Body.String())}
  var results []personalResult
  if err:=json.Unmarshal(w.Body.Bytes(),&results);err!=nil{t.Fatal(err)}
  if len(results)!=1{t.Fatalf("expected 1 result: %+v",results)}
  return results[0]
 }
 rss:=get(true,true)
 if len(rss.Articles)!=1||rss.Articles[0].Title!="Real story"{t.Fatalf("RSS headlines not verified: %+v",rss.Articles)}
 if len(rss.Filtered)!=1||rss.Filtered[0].Reason!="homepage"{t.Fatalf("RSS homepage not filtered: %+v",rss.Filtered)}
 if calls.Load()!=2{t.Fatalf("first run classification calls=%d",calls.Load())}
 full:=get(false,true)
 if len(full.Articles)!=1||len(full.Filtered)!=1{t.Fatalf("normal mode changed article gate: %+v",full)}
 if calls.Load()!=2{t.Fatalf("second mode failed to share cache; calls=%d",calls.Load())}
 bypass:=get(true,false)
 if len(bypass.Articles)!=1||len(bypass.Filtered)!=1{t.Fatalf("bypassed cache changed eligibility: %+v",bypass)}
 if calls.Load()!=4{t.Fatalf("opted-out RSS mode did not bypass cache; calls=%d",calls.Load())}
 again:=get(false,true)
 if len(again.Articles)!=1||len(again.Filtered)!=1||calls.Load()!=4{t.Fatalf("cache bypass invalidated other cached entries; calls=%d",calls.Load())}
}
