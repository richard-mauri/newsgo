package dashboard

import (
 "bytes"
 "fmt"
 "log"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

func TestNewsGoSharedGateAPIAndIdentityMatrix(t *testing.T) {
 token:=strings.Repeat("a",64)
 expired:=false
 upstreamCode:=200
 role:="guest"
 username:="guest-001"
 verifiedRequests:=0
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  verifiedRequests++
  if r.URL.Path!="/api/session/profile" || r.Header.Get("Authorization")!="Bearer "+token ||
   r.Header.Get("Origin")!=newsgoOrigin {t.Error("verification did not use fixed trusted Homepage profile contract")}
  if upstreamCode!=200{w.WriteHeader(upstreamCode);return}
  expires:=time.Now().Add(time.Hour)
  if expired {expires=time.Now().Add(-time.Minute)}
  fmt.Fprintf(w,`{"user":%q,"role":%q,"expiresAt":%q}`,username,role,expires.UTC().Format(time.RFC3339))
 }))
 defer upstream.Close()
 gate:=&NewsGoAuthGate{ProfileURL:upstream.URL+"/api/session/profile",Client:upstream.Client()}
 served:=0
 next:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){served++;w.WriteHeader(200)})
 app:=gate.Protect(next)
 call:=func(path,authHeader string,spoofedRole string)int{
  req:=httptest.NewRequest("GET","http://localhost:8084"+path,nil)
  if authHeader!=""{req.Header.Set("Authorization",authHeader)}
  req.Header.Set("X-NewsGo-Role",spoofedRole)
  rec:=httptest.NewRecorder()
  app.ServeHTTP(rec,req)
  if strings.HasPrefix(path,"/api")&&rec.Header().Get("Cache-Control")!="no-store"{t.Errorf("API response must not be cached: %q",path)}
  return rec.Code
 }
 for _,path:=range []string{"/api","/api/personal","/api/discover","/api/search-feeds","/api/diagnostic/feed-trace","/api/admin/settings","/api/admin/classification-diagnostic","/api/activity/ping","/api/not-yet-defined"} {
  if got:=call(path,"","admin");got!=401{t.Errorf("anonymous %s got %d",path,got)}
  if got:=call(path,"Bearer invalid","admin");got!=401{t.Errorf("malformed bearer %s got %d",path,got)}
  if got:=call(path,"Bearer "+token,"guest");got!=200{t.Errorf("verified Guest %s got %d",path,got)}
 }
 role="user";username="reader"
 if got:=call("/api/personal","Bearer "+token,"guest");got!=200{t.Fatalf("verified personal user got %d",got)}
 role="admin";username="operator"
 if got:=call("/api/admin/settings","Bearer "+token,"guest");got!=200{t.Fatalf("verified admin got %d",got)}
 role="untrusted"
 if got:=call("/api/personal","Bearer "+token,"admin");got!=401{t.Fatalf("unknown role accepted: %d",got)}
 role="user";username=""
 if got:=call("/api/personal","Bearer "+token,"admin");got!=401{t.Fatalf("empty user accepted: %d",got)}
 username="reader";expired=true
 if got:=call("/api/personal","Bearer "+token,"admin");got!=401{t.Fatalf("expired bearer accepted: %d",got)}
 expired=false;upstreamCode=401
 if got:=call("/api/personal","Bearer "+token,"admin");got!=401{t.Fatalf("revoked bearer accepted: %d",got)}
 upstreamCode=503
 if got:=call("/api/personal","Bearer "+token,"admin");got!=503{t.Fatalf("upstream outage should fail closed: %d",got)}
 upstreamCode=200
 for _,path:=range []string{"/","/personal.js","/suite-ui/banner.css","/health","/activity.js","/tanner.jpg"}{
  if got:=call(path,"","");got!=200{t.Errorf("public bootstrap/health %s inaccessible: %d",path,got)}
 }
 if verifiedRequests==0||served==0{t.Fatal("authorization matrix did not exercise verifier and handler")}
}

func TestNewsGoAuthGateRejectsNetworkFailure(t *testing.T){
 gate:=&NewsGoAuthGate{ProfileURL:"http://127.0.0.1:1/api/session/profile",Client:&http.Client{Timeout:time.Second}}
 handler:=gate.Protect(http.HandlerFunc(func(http.ResponseWriter,*http.Request){t.Fatal("unverified caller reached API")}))
 req:=httptest.NewRequest("GET","http://localhost:8084/api/personal",nil)
 req.Header.Set("Authorization","Bearer "+strings.Repeat("b",64))
 rec:=httptest.NewRecorder();handler.ServeHTTP(rec,req)
 if rec.Code!=503{t.Fatalf("unexpected network error status %d",rec.Code)}
}

func TestThumbnailBoundaryTracesBeforeAndInsideHandlerWithoutSecrets(t *testing.T) {
 token:=strings.Repeat("c",64)
 profile:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if got:=r.Header.Get("Authorization");got!="Bearer "+token {t.Error("profile bearer mismatch")}
  fmt.Fprintf(w,`{"user":"private-account","role":"user","expiresAt":%q}`,time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
 }))
 defer profile.Close()
 gate:=&NewsGoAuthGate{ProfileURL:profile.URL,Client:profile.Client()}
 app:=gate.Protect(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("X-NewsGo-Thumbnail-Proxy","v1")
  w.Header().Set("X-Thumbnail-Outcome","queue_timeout")
  w.Header().Set("X-Thumbnail-Queue-Ms","123")
  w.WriteHeader(http.StatusGatewayTimeout)
 }))
 var output bytes.Buffer
 original:=log.Writer()
 log.SetOutput(&output)
 defer log.SetOutput(original)
 request:=func(auth string) *httptest.ResponseRecorder {
  r:=httptest.NewRequest("GET","/api/thumbnail?url=https%3A%2F%2Fexample.com%2Fprivate",nil)
  if auth!="" {r.Header.Set("Authorization",auth)}
  r.Header.Set("X-NewsGo-Request-ID","untrusted-header")
  w:=httptest.NewRecorder()
  app.ServeHTTP(w,r)
  return w
 }
 unauth:=request("")
 if unauth.Code!=http.StatusUnauthorized{t.Fatalf("unauthorized returned %d",unauth.Code)}
 if unauth.Header().Get("X-NewsGo-Thumbnail-Proxy")!=""{t.Fatal("outer gate produced handler marker")}
 verified:=request("Bearer "+token)
 if verified.Code!=http.StatusGatewayTimeout{t.Fatalf("handler returned %d",verified.Code)}
 if verified.Header().Get("X-NewsGo-Thumbnail-Proxy")!="v1"{t.Fatal("handler marker lost")}
 id1,id2:=unauth.Header().Get("X-NewsGo-Request-ID"),verified.Header().Get("X-NewsGo-Request-ID")
 if id1==""||id2==""||id1==id2||id1=="untrusted-header"||id2=="untrusted-header"{t.Fatalf("bad IDs %q %q",id1,id2)}
 logs:=output.String()
 for _,want:=range []string{"THUMBNAIL boundary id="+id1+" status=401 handler_marked=false outcome=before_handler","THUMBNAIL boundary id="+id2+" status=504 handler_marked=true outcome=queue_timeout","queue_ms=123"} {
  if !strings.Contains(logs,want){t.Errorf("missing trace fragment %q in %s",want,logs)}
 }
 for _,secret:=range []string{token,"private-account","example.com","private","untrusted-header"}{
  if strings.Contains(logs,secret){t.Errorf("log exposed private content: %q",secret)}
 }
}

func TestThumbnailBoundarySuccessesAreAggregated(t *testing.T) {
 var output bytes.Buffer
 original:=log.Writer()
 log.SetOutput(&output)
 defer log.SetOutput(original)
 handler:=traceThumbnailBoundary(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("X-NewsGo-Thumbnail-Proxy","v1")
  w.WriteHeader(http.StatusOK)
 }))
 start:=thumbnailCompletedSuccesses.Load()
 for i:=0;i<100;i++{
  w:=httptest.NewRecorder()
  handler.ServeHTTP(w,httptest.NewRequest("GET","/api/thumbnail?url=https://example.com",nil))
  if w.Code!=200{t.Fatal("expected success")}
 }
 logs:=output.String()
 if strings.Contains(logs,"boundary id="){t.Fatal("successful thumbnail logged individually")}
 if !strings.Contains(logs,"THUMBNAIL boundary_success_total="){t.Fatalf("expected periodic success aggregate, start=%d: %s",start,logs)}
}

func TestThumbnailHandlerAndBoundaryShareServerID(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN","")
 token:=strings.Repeat("d",64)
 profile:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  fmt.Fprintf(w,`{"user":"test","role":"user","expiresAt":%q}`,time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
 }))
 defer profile.Close()
 gate:=&NewsGoAuthGate{ProfileURL:profile.URL,Client:profile.Client()}
 app:=gate.Protect(http.HandlerFunc((&State{reader:"https://scrapergo.themauris.org"}).serveThumbnail))
 var output bytes.Buffer
 old:=log.Writer();log.SetOutput(&output);defer log.SetOutput(old)
 req:=httptest.NewRequest("GET","/api/thumbnail?url=https%3A%2F%2Fexample.com%2Fprivate-story",nil)
 req.Header.Set("Authorization","Bearer "+token)
 rec:=httptest.NewRecorder()
 app.ServeHTTP(rec,req)
 if rec.Code!=http.StatusBadGateway{t.Fatalf("got %d",rec.Code)}
 id:=rec.Header().Get("X-NewsGo-Request-ID")
 if id==""{t.Fatal("missing request ID")}
 logText:=output.String()
 for _,marker:=range []string{"THUMBNAIL handler id="+id+" status=502","THUMBNAIL boundary id="+id+" status=502"} {
  if !strings.Contains(logText,marker){t.Errorf("missing correlated trace: %s",marker)}
 }
 for _,secret:=range []string{token,"example.com","private-story"}{
  if strings.Contains(logText,secret){t.Fatal("sensitive request data in log")}
 }
}
