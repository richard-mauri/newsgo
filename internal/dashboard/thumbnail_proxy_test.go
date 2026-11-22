package dashboard

import (
 "net/http"
 "net/http/httptest"
 "net/url"
 "strings"
 "testing"
 "time"
 "strconv"
)

const testThumbnailToken = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
const testThumbnailArticle = "https://example.com/news/one"

func TestThumbnailProxyRequiresSharedCredentialAndValidRequest(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN","")
 state:=&State{reader:"https://scrapergo.themauris.org"}
 for _,tc:=range []struct{method,article,width string; want int}{
  {"GET",testThumbnailArticle,"256",http.StatusBadGateway},
  {"POST",testThumbnailArticle,"256",http.StatusMethodNotAllowed},
  {"GET","http://127.0.0.1/private","256",http.StatusBadRequest},
  {"GET","ftp://example.com/a","256",http.StatusBadRequest},
  {"GET",testThumbnailArticle,"999",http.StatusBadRequest},
 }{
  path:="/api/thumbnail?url="+url.QueryEscape(tc.article)+"&width="+tc.width
  w:=httptest.NewRecorder()
  state.serveThumbnail(w,httptest.NewRequest(tc.method,path,nil))
  if w.Code!=tc.want{t.Errorf("%s %q width=%s: HTTP %d, want %d",tc.method,tc.article,tc.width,w.Code,tc.want)}
 }
}

func TestThumbnailProxyPreservesImageAndMachineCredential(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 payload:=[]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0}
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/thumbnail"{t.Errorf("unexpected path %q",r.URL.Path)}
  if r.Header.Get("Authorization")!="Service "+testThumbnailToken{t.Error("missing thumbnail service authorization")}
  if r.URL.Query().Get("url")!=testThumbnailArticle||r.URL.Query().Get("width")!="256"{t.Error("query mismatch")}
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write(payload)
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 w:=httptest.NewRecorder()
 state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle)+"&width=256",nil))
 if w.Code!=http.StatusOK{t.Fatalf("HTTP %d: %s",w.Code,w.Body.String())}
 if w.Header().Get("Content-Type")!="image/png"{t.Errorf("unexpected MIME: %q",w.Header().Get("Content-Type"))}
 if w.Header().Get("Cache-Control")!="no-store"{t.Error("image response must not be stored")}
 if w.Body.String()!=string(payload){t.Error("image bytes changed")}
 if strings.Contains(w.Body.String(),testThumbnailToken){t.Error("service credential exposed")}
}

func TestThumbnailProxyRejectsNonimageAndRedirect(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 for _,tc:=range []struct{status int;mime string}{
  {http.StatusFound,"image/png"},
  {http.StatusOK,"text/html"},
  {http.StatusUnauthorized,"text/plain"},
 }{
  upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
   w.Header().Set("Content-Type",tc.mime)
   if tc.status==http.StatusFound{w.Header().Set("Location","https://example.org/should-not-follow")}
   w.WriteHeader(tc.status)
   _,_=w.Write([]byte("not image"))
  }))
  state:=&State{reader:upstream.URL}
  w:=httptest.NewRecorder()
  state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
  if w.Code!=http.StatusBadGateway{t.Errorf("upstream %d %s yielded HTTP %d",tc.status,tc.mime,w.Code)}
  upstream.Close()
 }
}

func TestThumbnailURLUsesProtectedNewsGoRoute(t *testing.T) {
 got:=scraperThumbnailURL("https://scrapergo.themauris.org",testThumbnailArticle)
 if !strings.HasPrefix(got,"/api/thumbnail?"){t.Fatalf("thumbnail URL is not local: %q",got)}
 parsed,err:=url.Parse(got)
 if err!=nil||parsed.Query().Get("url")!=testThumbnailArticle||parsed.Query().Get("width")!="256"{t.Fatalf("bad URL: %q",got)}
 if scraperThumbnailURL("https://scrapergo.themauris.org","file:///etc/passwd")!=""{t.Error("unsafe URL accepted")}
}

func TestThumbnailLoopbackBackendIsolation(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 const publicReader="https://scrapergo.themauris.org"
 payload:=[]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0}
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/thumbnail" {t.Errorf("unexpected service path: %q",r.URL.Path)}
  if r.Header.Get("Authorization")!="Service "+testThumbnailToken {t.Error("shared service credential absent")}
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write(payload)
 }))
 defer upstream.Close()
 state:=&State{reader:publicReader}
 if err:=state.SetThumbnailBackend(upstream.URL);err!=nil {t.Fatal(err)}
 if state.reader!=publicReader {t.Fatalf("public reader changed: %q",state.reader)}
 w:=httptest.NewRecorder()
 state.serveThumbnail(w,httptest.NewRequest(http.MethodGet,"/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
 if w.Code!=http.StatusOK {t.Fatalf("loopback status %d: %s",w.Code,w.Body.String())}
 if !strings.HasPrefix(scraperThumbnailURL(state.reader,testThumbnailArticle),"/api/thumbnail?") {t.Fatal("public reader thumbnail route changed")}
}

func TestThumbnailLoopbackRejectsUnsafeDestinations(t *testing.T) {
 state:=&State{reader:"https://scrapergo.themauris.org"}
 for _,raw:=range []string{
  "https://127.0.0.1:8083","http://localhost:8083",
  "http://192.168.1.1:8083","http://127.0.0.1",
  "http://127.0.0.1:8083/other","http://127.0.0.1:8083?x=y",
  "http://user:password@127.0.0.1:8083",
  "http://127.0.0.1:8083#fragment",
  "http://scrapergo.themauris.org",
 } {
  if err:=state.SetThumbnailBackend(raw);err==nil {t.Errorf("accepted unsafe thumbnail backend %q",raw)}
  if state.thumbnailBackend!="" {t.Errorf("invalid backend retained after %q",raw)}
 }
 if err:=state.SetThumbnailBackend("http://127.0.0.1:8083");err!=nil {t.Fatal(err)}
 if err:=state.SetThumbnailBackend("");err!=nil || state.thumbnailBackend!="" {t.Fatal("disable backend failed")}
}

func TestThumbnailQueueAdmissionTimeoutIndependentFromUpstream(t *testing.T) {
 previous:=currentThumbnailSettings()
 activeThumbnailSettings.Store(ThumbnailSettings{Concurrency:1,QueueWaitSeconds:1,UpstreamTimeoutSeconds:2})
 defer activeThumbnailSettings.Store(previous)
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 entered:=make(chan struct{},1)
 release:=make(chan struct{})
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  entered<-struct{}{}
  <-release
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write([]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0})
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 first:=make(chan *httptest.ResponseRecorder,1)
 go func(){
  w:=httptest.NewRecorder()
  state.serveThumbnail(w,httptest.NewRequest(http.MethodGet,"/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
  first<-w
 }()
 <-entered
 second:=httptest.NewRecorder()
 started:=time.Now()
 state.serveThumbnail(second,httptest.NewRequest(http.MethodGet,"/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
 if second.Code!=http.StatusGatewayTimeout || second.Header().Get("X-Thumbnail-Outcome")!="queue_timeout" {
  t.Fatalf("queue timeout: code=%d reason=%q",second.Code,second.Header().Get("X-Thumbnail-Outcome"))
 }
 if time.Since(started)<900*time.Millisecond {t.Fatal("queue timeout returned before admission budget")}
 close(release)
 select {
 case w:=<-first:
  if w.Code!=http.StatusOK || w.Header().Get("X-Thumbnail-Outcome")!="loaded" {t.Fatalf("first response %d",w.Code)}
 case <-time.After(3*time.Second):t.Fatal("first request did not complete")
 }
 active,_:=thumbnailActivePeak()
 if active!=0 {t.Fatalf("leaked thumbnail slot: %d",active)}
}

func TestThumbnailUpstreamDeadlineAndSafeDiagnosticHeaders(t *testing.T) {
 previous:=currentThumbnailSettings()
 activeThumbnailSettings.Store(ThumbnailSettings{Concurrency:6,QueueWaitSeconds:1,UpstreamTimeoutSeconds:2})
 defer activeThumbnailSettings.Store(previous)
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  <-r.Context().Done()
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 w:=httptest.NewRecorder()
 state.serveThumbnail(w,httptest.NewRequest(http.MethodGet,"/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
 if w.Code!=http.StatusGatewayTimeout || w.Header().Get("X-Thumbnail-Outcome")!="upstream_timeout" {
  t.Fatalf("upstream timeout %d reason %q",w.Code,w.Header().Get("X-Thumbnail-Outcome"))
 }
 for _,name:=range []string{"X-Thumbnail-Queue-Ms","X-Thumbnail-Upstream-Ms"} {
  if _,err:=strconv.Atoi(w.Header().Get(name));err!=nil {t.Errorf("%s must be numeric: %v",name,err)}
 }
 if strings.Contains(w.Body.String(),testThumbnailToken) || strings.Contains(w.Body.String(),testThumbnailArticle) {t.Fatal("sensitive request data in response")}
}

func TestThumbnailProxyResponseOriginMarkerOnSuccessAndFailures(t *testing.T){
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 validPNG:=[]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0}
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Content-Type","image/png");_,_=w.Write(validPNG)
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 for _,tc:=range []struct{name,method,path string;status int;outcome string}{
  {"success","GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),http.StatusOK,"loaded"},
  {"invalid-article","GET","/api/thumbnail?url="+url.QueryEscape("http://127.0.0.1/"),http.StatusBadRequest,""},
  {"bad-method","POST","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),http.StatusMethodNotAllowed,""},
 }{
  t.Run(tc.name,func(t *testing.T){
   w:=httptest.NewRecorder()
   state.serveThumbnail(w,httptest.NewRequest(tc.method,tc.path,nil))
   if w.Code!=tc.status{t.Fatalf("got %d want %d",w.Code,tc.status)}
   if v:=w.Header().Get("X-NewsGo-Thumbnail-Proxy");v!="v1"{t.Fatalf("missing fixed origin marker %q",v)}
   if w.Header().Get("X-Thumbnail-Outcome")!=tc.outcome{t.Fatal("unexpected classified outcome")}
   for _,name:=range []string{"X-NewsGo-Thumbnail-Proxy","X-Thumbnail-Outcome","X-Thumbnail-Queue-Ms"}{
    if strings.Contains(w.Header().Get(name),testThumbnailToken)||strings.Contains(w.Header().Get(name),testThumbnailArticle){t.Fatalf("sensitive data leaked in %s",name)}
   }
  })
 }
}

func TestThumbnailProxyQueueTimeoutIsMarked(t *testing.T){
 previous:=currentThumbnailSettings()
 activeThumbnailSettings.Store(ThumbnailSettings{Concurrency:1,QueueWaitSeconds:1,UpstreamTimeoutSeconds:2})
 defer activeThumbnailSettings.Store(previous)
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 entered:=make(chan struct{},1);release:=make(chan struct{})
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  entered<-struct{}{};<-release
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write([]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0})
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 first:=make(chan struct{})
 go func(){defer close(first);w:=httptest.NewRecorder();state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))}()
 <-entered
 w:=httptest.NewRecorder();state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
 if w.Code!=http.StatusGatewayTimeout || w.Header().Get("X-Thumbnail-Outcome")!="queue_timeout" || w.Header().Get("X-NewsGo-Thumbnail-Proxy")!="v1" {
  t.Fatalf("queue 504 must be marked: %d, %q, %q",w.Code,w.Header().Get("X-Thumbnail-Outcome"),w.Header().Get("X-NewsGo-Thumbnail-Proxy"))
 }
 close(release)
 select{case <-first:case <-time.After(3*time.Second):t.Fatal("upstream request hung")}
}
