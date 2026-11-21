package dashboard

import (
 "fmt"
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
