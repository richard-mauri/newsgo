package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "sync/atomic"
 "strings"
 "testing"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestExcludedDomainNeverReclassified(t *testing.T) {
 previous:=currentArticleSettings()
 cfg:=defaultArticleSettings()
 cfg.ClassificationCacheMinutes=720
 activeArticleSettings.Store(cfg)
 defer activeArticleSettings.Store(previous)
 var calls atomic.Int64
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls.Add(1)
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"unknown"}`))
 }))
 defer upstream.Close()
 articles:=[]feeds.Article{
  {URL:"https://www.investors.com/news/a"},
  {URL:"https://sub.investors.com/news/b"},
  {URL:"https://investors.com/news/c"},
  {URL:"https://notinvestors.com/news/d"},
 }
 for i:=0;i<2;i++ {
  stats:=&classificationCacheStats{}
  outcomes:=classifyFeedArticlesWithExclusionsStats(context.Background(),upstream.URL,articles,0,true,[]string{"investors.com"},stats)
  for k:=0;k<3;k++ {if outcomes[k]!="publisher exclusion"{t.Fatalf("excluded outcome at %d: %q",k,outcomes[k])}}
  if outcomes[3]!="unknown"{t.Fatalf("unexcluded article incorrectly skipped: %q",outcomes[3])}
  if stats.misses.Load()!=1||stats.cold.Load()!=1||stats.unknown.Load()!=1||stats.hits.Load()!=0{
   t.Fatalf("skipped domains counted as classification requests: misses=%d cold=%d unknown=%d hits=%d",stats.misses.Load(),stats.cold.Load(),stats.unknown.Load(),stats.hits.Load())
  }
 }
 if calls.Load()!=2 {t.Fatalf("excluded domains were sent to ScraperGo: calls=%d; expected only two nonexcluded requests",calls.Load())}
}

func TestPersonalDomainExclusionHostBoundary(t *testing.T){
 if !validPersonalExcludedDomains([]string{"investors.com","bom.gov.au","youtube.com"}){t.Fatal("valid normalized excluded domains rejected")}
 for _,domains:=range [][]string{
  {"http://investors.com"},{"investors.com/path"},{"investors.com:443"},
  {"Investors.com"},{"localhost"},{"investors.com","bad space.org"},
 } {
  if validPersonalExcludedDomains(domains){t.Fatalf("invalid excluded domains allowed: %v",domains)}
 }
 for _,raw:=range []string{"https://investors.com/story","https://www.investors.com/story","https://sub.investors.com/story"} {
  if !excludedArticleURL(raw,[]string{"investors.com"}){t.Fatalf("excluded hostname not matched: %s",raw)}
 }
 for _,raw:=range []string{"https://notinvestors.com/story","https://investors.com.evil.org/story","https://example.org/story"} {
  if excludedArticleURL(raw,[]string{"investors.com"}){t.Fatalf("unrelated hostname incorrectly excluded: %s",raw)}
 }
}

func TestPersonalReaderSendsExcludedDomainsBeforeClassification(t *testing.T) {
 if !strings.Contains(personalJS, "excludedDomains:effectiveExcludedDomains()") {
  t.Fatal("NewsGo browser must send its normalized exclusion list to /api/personal")
 }
 if !strings.Contains(personalJS, "const res=await fetch('/api/personal'") {
  t.Fatal("missing bounded personal feed read request")
 }
}
