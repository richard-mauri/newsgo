package dashboard

import (
	"encoding/json"
 "context"
 "io"
 "net/http"
 "sync/atomic"
 "time"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseFeedlyResults(t *testing.T) {
	body := []byte(`{"results":[{"title":"Sailing","feedId":"feed/https://example.org/rss","website":"https://example.org"},{"title":"Duplicate","id":"feed/https://example.org/rss"},{"title":"Private","feedId":"feed/http://127.0.0.1/private"}]}`)
	got, err := feedlyResults(body)
	if err != nil || len(got) != 1 || got[0].Name != "Sailing" {
		t.Fatalf("unexpected %#v %v", got, err)
	}
}
func TestSearchRejectsBadQuery(t *testing.T) {
	s := &State{}
	w := httptest.NewRecorder()
	s.serveFeedSearch(w, httptest.NewRequest("GET", "/api/search-feeds?q=x", nil))
	if w.Code != 400 {
		t.Fatalf("HTTP %d", w.Code)
	}
	var r feedSearchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil || !strings.Contains(r.Error, "topic") {
		t.Fatalf("%v %+v", err, r)
	}
}

type directoryTestTransport func(*http.Request)(*http.Response,error)
func (fn directoryTestTransport)RoundTrip(req *http.Request)(*http.Response,error){return fn(req)}
func directoryTestHTTPClient(status int,body string)*http.Client{
 return &http.Client{Transport:directoryTestTransport(func(req *http.Request)(*http.Response,error){
  return &http.Response{StatusCode:status,Status:http.StatusText(status),Header:http.Header{"Content-Type":[]string{"application/rss+xml"}},Body:io.NopCloser(strings.NewReader(body)),Request:req},nil
 })}
}
func TestDirectoryFeedPrevalidation(t *testing.T){
 validRSS:=`<rss version="2.0"><channel><item><title>Story</title><link>https://example.org/story</link></item></channel></rss>`
 for _,tc:=range []struct{name string;status int;body string;wantState string;wantReason string}{
  {"working",200,validRSS,"available",""},
  {"obsolete400",400,"","unavailable","HTTP 400"},
  {"gone404",404,"","unavailable","HTTP 404"},
  {"publisherOutage503",503,"","unverified","HTTP 503"},
  {"rateLimited429",429,"","unverified","HTTP 429"},
  {"badMarkup",200,"<html>not a feed</html>","unavailable","Invalid or oversized feed response"},
  {"emptyFeed",200,`<rss><channel></channel></rss>`,"unavailable","No readable articles in this feed"},
 } {
  t.Run(tc.name,func(t *testing.T){
   item:=verifyDirectoryFeed(context.Background(),directoryTestHTTPClient(tc.status,tc.body),feedSearchItem{Name:"News",URL:"https://example.org/rss"})
   if item.Status!=tc.wantState||item.Reason!=tc.wantReason {t.Fatalf("status=%q reason=%q, wanted %q / %q",item.Status,item.Reason,tc.wantState,tc.wantReason)}
  })
 }
}
func TestDirectoryCandidateCheckIsBoundedAndIndependent(t *testing.T){
 items:=make([]feedSearchItem,12)
 for i:=range items{items[i]=feedSearchItem{Name:"Candidate",URL:"https://example.org/rss"}}
 var current,peak atomic.Int32
 out:=verifyDirectoryCandidates(context.Background(),items,func(ctx context.Context,item feedSearchItem)feedSearchItem{
  active:=current.Add(1)
  for {old:=peak.Load();if active<=old||peak.CompareAndSwap(old,active){break}}
  time.Sleep(4*time.Millisecond)
  current.Add(-1)
  item.Status="available"
  return item
 })
 if peak.Load()>4||peak.Load()<2 {t.Fatalf("unexpected concurrent validations %d",peak.Load())}
 if len(out)!=12 {t.Fatalf("got %d",len(out))}
 for _,item:=range out {if item.Status!="available"{t.Fatalf("unexpected result %+v",item)}}
}
