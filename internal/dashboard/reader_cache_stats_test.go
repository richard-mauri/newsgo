package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"
)

func TestReaderClassificationCacheStats(t *testing.T){
 previous:=currentArticleSettings()
 cfg:=defaultArticleSettings()
 cfg.ClassificationCacheMinutes=1440
 activeArticleSettings.Store(cfg)
 defer activeArticleSettings.Store(previous)
 clearNewsClassificationCache()
 defer clearNewsClassificationCache()
 calls:=0
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls++
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer upstream.Close()
 url:="https://example.org/cache-stats-test"
 stats:=&classificationCacheStats{}
 if got:=newsClassifyWithCacheStats(context.Background(),upstream.URL,url,true,stats);got!="article"{t.Fatal(got)}
 if stats.misses.Load()!=1||stats.hits.Load()!=0||stats.bypasses.Load()!=0{t.Fatal("cold miss not counted")}
 if got:=newsClassifyWithCacheStats(context.Background(),upstream.URL,url,true,stats);got!="article"{t.Fatal(got)}
 if stats.hits.Load()!=1||calls!=1{t.Fatalf("hit not counted or sent upstream: hits=%d calls=%d",stats.hits.Load(),calls)}
 if got:=newsClassifyWithCacheStats(context.Background(),upstream.URL,url,false,stats);got!="article"{t.Fatal(got)}
 if stats.bypasses.Load()!=1||calls!=2{t.Fatalf("bypass not counted or sent upstream: bypass=%d calls=%d",stats.bypasses.Load(),calls)}
 if got:=newsClassifyWithCacheStats(context.Background(),upstream.URL,url,true,stats);got!="article"{t.Fatal(got)}
 if stats.hits.Load()!=2||calls!=2{t.Fatalf("bypass modified cache: hits=%d calls=%d",stats.hits.Load(),calls)}
 _=time.Second
}
