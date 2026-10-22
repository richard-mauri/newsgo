package dashboard

import(
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
 "github.com/richard-mauri/newsgo/internal/feeds"
)
func TestFeedTraceMissingOutsideCapAndAge(t *testing.T){
 now:=time.Date(2026,10,21,17,0,0,0,time.UTC)
 items:=[]feeds.Article{
  {Title:"Older",URL:"https://www.npr.org/old",Published:now.AddDate(0,0,-40)},
  {Title:"West Bank",URL:"https://www.npr.org/2026/10/21/nx-s1-6010460/brad-lander-detained-west-bank-israel-palestinians",Published:now},
 }
 target:=items[1].URL+"?utm_source=whatever"
 if v:=traceArticleInFeed(items,target,1,30,now);v.Status!="outside-processing-cap"||v.MatchIndex!=1{t.Fatalf("processing cap: %+v",v)}
 if v:=traceArticleInFeed(items,target,5,30,now);v.Status!="eligible-for-classification" {t.Fatalf("eligibility: %+v",v)}
 if v:=traceArticleInFeed(items,items[0].URL,5,30,now);v.Status!="age-filtered"{t.Fatalf("age: %+v",v)}
 if v:=traceArticleInFeed(items,"https://www.npr.org/missing",5,30,now);v.Status!="not-in-parsed-feed"{t.Fatalf("absent: %+v",v)}
}
func TestFeedTraceRejectsRemotePeer(t *testing.T){
 s:=&State{}
 r:=httptest.NewRequest("POST","/api/diagnostic/feed-trace",strings.NewReader(`{}`))
 r.RemoteAddr="203.0.113.12:5000"
 w:=httptest.NewRecorder()
 s.serveFeedTrace(w,r)
 if w.Code!=http.StatusForbidden{t.Fatalf("status %d",w.Code)}
}
