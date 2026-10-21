package dashboard

import (
 "bytes"
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestAgeFilteringSkipsOldBeforeClassification(t *testing.T) {
 now := time.Now()
 articles := []feeds.Article{
  {Title:"Old daily news",URL:"https://example.org/old",Published:now.Add(-20*24*time.Hour)},
  {Title:"Undated story",URL:"https://example.org/undated"},
 }
 decisions := classifyFeedArticlesWithAge(context.Background(),"http://127.0.0.1:1",articles,3)
 if decisions[0]!="too-old" {t.Fatalf("old story classified as %q, want too-old",decisions[0])}
 if decisions[1]=="too-old" {t.Fatal("undated story filtered by age")}
}

func TestPersonalAPIAgeFilterAndAudit(t *testing.T) {
 now:=time.Now().UTC()
 source:=feeds.Source{Name:"Age Test",Category:"Other",URL:"https://example.org/age-feed"}
 previous:=currentArticleSettings()
 defer activeArticleSettings.Store(previous)
 activeArticleSettings.Store(defaultArticleSettings())
 personalCache.Lock()
 old, exists:=personalCache.entries[source.URL]
 personalCache.entries[source.URL]=cacheEntry{articles:[]feeds.Article{
  {Title:"Old daily news",URL:"https://example.org/old",Published:now.Add(-20*24*time.Hour)},
  {Title:"New story",URL:"https://example.org/new",Published:now.Add(-time.Hour)},
 },at:time.Now()}
 personalCache.Unlock()
 defer func(){personalCache.Lock();if exists {personalCache.entries[source.URL]=old}else{delete(personalCache.entries,source.URL)};personalCache.Unlock()}()
 // The classifier endpoint deliberately marks new articles as valid; old
 // articles must never be sent to it.
 calls:=0
 reader:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls++
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer reader.Close()
 state:=&State{reader:reader.URL}
 body,_:=json.Marshal(personalRequest{Feeds:[]feeds.Source{source},MaxArticleAgeDays:3})
 req:=httptest.NewRequest(http.MethodPost,"/api/personal",bytes.NewReader(body))
 res:=httptest.NewRecorder()
 state.servePersonal(res,req)
 if res.Code!=http.StatusOK {t.Fatalf("HTTP %d: %s",res.Code,res.Body.String())}
 var result []personalResult
 if err:=json.Unmarshal(res.Body.Bytes(),&result);err!=nil {t.Fatal(err)}
 if len(result)!=1 {t.Fatalf("got %d groups",len(result))}
 if len(result[0].Articles)!=1||result[0].Articles[0].Title!="New story" {t.Fatalf("unexpected visible: %+v",result[0].Articles)}
 if len(result[0].Filtered)!=1||result[0].Filtered[0].Reason!="too-old"||result[0].Filtered[0].Title!="Old daily news" {t.Fatalf("unexpected filtered: %+v",result[0].Filtered)}
 if calls!=1 {t.Fatalf("classification requests=%d, want 1",calls)}
}
