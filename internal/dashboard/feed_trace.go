package dashboard

import (
 "encoding/json"
 "net"
 "net/http"
 "net/url"
 "strings"
 "time"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

type feedTraceRequest struct {
 ArticleURL string `json:"articleUrl"`
 Feeds []feeds.Source `json:"feeds"`
 MaxArticleAgeDays int `json:"maxArticleAgeDays"`
}
type feedTraceResult struct {
 Feed string `json:"feed"`
 FeedURL string `json:"feedUrl"`
 Error string `json:"error,omitempty"`
 Entries int `json:"entries"`
 MatchIndex int `json:"matchIndex"`
 MatchedURL string `json:"matchedUrl,omitempty"`
 Title string `json:"title,omitempty"`
 Published string `json:"published,omitempty"`
 Cutoff int `json:"processingCap"`
 Status string `json:"status"`
 Classification string `json:"classification,omitempty"`
 CacheAgeLimitDays int `json:"ageLimitDays"`
}
func comparableArticleURL(raw string) string {
 u,err:=url.Parse(strings.TrimSpace(raw));if err!=nil{return ""}
 u.Fragment=""
 u.RawQuery=""
 return strings.TrimSuffix(strings.ToLower(u.Hostname())+u.EscapedPath(),"/")
}
func traceArticleInFeed(articles []feeds.Article, target string, cap, ageDays int, now time.Time) feedTraceResult {
 result:=feedTraceResult{Entries:len(articles),MatchIndex:-1,Cutoff:cap,CacheAgeLimitDays:ageDays,Status:"not-in-parsed-feed"}
 key:=comparableArticleURL(target)
 for i,a:=range articles {
  if comparableArticleURL(a.URL)!=key {continue}
  result.MatchIndex=i;result.MatchedURL=a.URL;result.Title=a.Title
  if !a.Published.IsZero(){result.Published=a.Published.Format(time.RFC3339)}
  switch {
  case i>=cap:result.Status="outside-processing-cap"
  case articleTooOld(a.Published,ageDays,now):result.Status="age-filtered"
  default:result.Status="eligible-for-classification"
  }
  break
 }
 return result
}
// Loopback-only operational inspection of explicit feed subscriptions. Never
// exposes stored subscriptions, mutates settings or bypasses safeFeedClient.
func(s *State) serveFeedTrace(w http.ResponseWriter,r *http.Request) {
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 w.Header().Set("Cache-Control","no-store")
 if r.Method!=http.MethodPost {w.Header().Set("Allow","POST");http.Error(w,"POST required",405);return}
 host,_,err:=net.SplitHostPort(r.RemoteAddr)
 if err!=nil||net.ParseIP(host)==nil||!net.ParseIP(host).IsLoopback(){http.Error(w,"localhost only",403);return}
 r.Body=http.MaxBytesReader(w,r.Body,8192)
 var req feedTraceRequest
 dec:=json.NewDecoder(r.Body)
 if err:=dec.Decode(&req);err!=nil||!safeURL(req.ArticleURL)||len(req.Feeds)==0||len(req.Feeds)>8||!validArticleAgeDays(req.MaxArticleAgeDays) {http.Error(w,"invalid diagnostic request",400);return}
 results:=make([]feedTraceResult,0,len(req.Feeds))
 for _,src:=range req.Feeds {
  result:=feedTraceResult{Feed:src.Name,FeedURL:src.URL,MatchIndex:-1,Status:"invalid-feed"}
  if len(src.Name)==0||len(src.Name)>70||len(src.Category)>40||len(src.URL)>2048||!safeURL(src.URL) {results=append(results,result);continue}
  articles,msg:=personalFetch(r.Context(),src)
  cap:=currentArticleSettings().MaxArticlesPerFeed
  result=traceArticleInFeed(articles,req.ArticleURL,cap,req.MaxArticleAgeDays,time.Now())
  result.Feed=src.Name;result.FeedURL=src.URL;result.Error=msg
  if msg!="" {result.Status="feed-fetch-error"}
  if result.Status=="eligible-for-classification" {result.Classification=newsClassify(r.Context(),s.reader,result.MatchedURL);if result.Classification=="article" {result.Status="classified-article"} else {result.Status="classification-filtered"}}
  results=append(results,result)
 }
 w.Header().Set("X-Content-Type-Options","nosniff")
 _=json.NewEncoder(w).Encode(results)
}
