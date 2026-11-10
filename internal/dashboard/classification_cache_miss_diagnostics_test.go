package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "testing"
 "time"
)

func TestClassificationCacheRepeatedNoncacheableOutcomesAreExplained(t *testing.T) {
 previous:=currentArticleSettings()
 cfg:=defaultArticleSettings()
 cfg.ClassificationCacheMinutes=720
 activeArticleSettings.Store(cfg)
 defer activeArticleSettings.Store(previous)
 clearNewsClassificationCache()
 defer clearNewsClassificationCache()
 outcome:="unknown"
 calls:=0
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls++
  if outcome=="failure" {http.Error(w,"temporary",http.StatusBadGateway);return}
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"`+outcome+`"}`))
 }))
 defer server.Close()

 stats:=&classificationCacheStats{}
 url:="https://example.org/repeated-unknown"
 for i:=0;i<2;i++ {
  if got:=newsClassifyWithCacheStats(context.Background(),server.URL,url,true,stats);got!="unknown" {t.Fatalf("unknown result: %s",got)}
 }
 if stats.misses.Load()!=2||stats.cold.Load()!=2||stats.hits.Load()!=0||stats.unknown.Load()!=2||calls!=2 {
  t.Fatalf("unknown must be retried and reported: misses=%d cold=%d hits=%d unknown=%d calls=%d",stats.misses.Load(),stats.cold.Load(),stats.hits.Load(),stats.unknown.Load(),calls)
 }
 outcome="unavailable"
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,"https://example.org/unavailable",true,stats);got!="unavailable" {t.Fatal(got)}
 if stats.unavailable.Load()!=1 {t.Fatalf("unavailable not explained: %d",stats.unavailable.Load())}
 outcome="failure"
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,"https://example.org/failed",true,stats);got!="unknown" {t.Fatal(got)}
 if stats.failures.Load()!=1 {t.Fatalf("upstream failure not explained: %d",stats.failures.Load())}
 outcome="article"
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,"https://example.org/cacheable",true,stats);got!="article"{t.Fatal(got)}
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,"https://example.org/cacheable",true,stats);got!="article"{t.Fatal(got)}
 if stats.hits.Load()!=1||stats.expired.Load()!=0 {t.Fatalf("successful result not cached: hits=%d expired=%d",stats.hits.Load(),stats.expired.Load())}
 if stats.cold.Load()!=stats.misses.Load() {t.Fatal("cold/expired breakdown not reconciled")}
}

func TestClassificationCacheExpiredMissIsReported(t *testing.T) {
 previous:=currentArticleSettings()
 cfg:=defaultArticleSettings()
 cfg.ClassificationCacheMinutes=720
 activeArticleSettings.Store(cfg)
 defer activeArticleSettings.Store(previous)
 clearNewsClassificationCache()
 defer clearNewsClassificationCache()
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer server.Close()
 url:="https://example.org/expired-diagnostic"
 stats:=&classificationCacheStats{}
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,url,true,stats);got!="article"{t.Fatal(got)}
 key:=server.URL+"|"+normalizeClassificationURL(url)
 newsClassificationCache.Lock()
 entry:=newsClassificationCache.entries[key]
 entry.expiry=time.Now().Add(-time.Second)
 newsClassificationCache.entries[key]=entry
 newsClassificationCache.Unlock()
 if got:=newsClassifyWithCacheStats(context.Background(),server.URL,url,true,stats);got!="article"{t.Fatal(got)}
 if stats.cold.Load()!=1||stats.expired.Load()!=1||stats.misses.Load()!=2 {
  t.Fatalf("miss breakdown: cold=%d expired=%d total=%d",stats.cold.Load(),stats.expired.Load(),stats.misses.Load())
 }
}
