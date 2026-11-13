package dashboard

import (
 "context"
 "encoding/json"
 "os"
 "os/exec"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"

 "github.com/richard-mauri/newsgo/internal/feeds"
)

func TestActivityReporterConfig(t *testing.T){
 if p,err:=NewActivityReporter("","");err!=nil||p!=nil{t.Fatalf("disabled: %v",err)}
 for _,url:=range []string{"http://example.com/api/admin/activity/ingest","http://localhost:8788/not-ingest","https://user:password@example.com/api/admin/activity/ingest"}{
  if _,err:=NewActivityReporter(url,strings.Repeat("x",40));err==nil{t.Fatalf("unsafe collector accepted %s",url)}
 }
 if _,err:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest","short");err==nil{t.Fatal("weak secret accepted")}
 if p,err:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("x",40));err!=nil||p==nil{t.Fatalf("loopback config rejected: %v",err)}
}

func TestActivityPingAndRollingCounts(t *testing.T){
 a,err:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("x",40))
 if err!=nil{t.Fatal(err)}
 now:=time.Now().UTC().Truncate(time.Second)
 a.now=func()time.Time{return now}
 ping:=func(origin,id string)int{
  r:=httptest.NewRequest("POST","http://127.0.0.1:9081/api/activity/ping",strings.NewReader(`{"id":"`+id+`","visible":true,"interacted":true}`))
  r.Header.Set("Origin",origin);r.Header.Set("Content-Type","application/json")
  w:=httptest.NewRecorder();a.ping(w,r);return w.Code
 }
 id:=strings.Repeat("a",32)
 if status:=ping("http://other.example",id);status!=403{t.Fatalf("cross origin accepted: %d",status)}
 if status:=ping("http://127.0.0.1:9081",id);status!=204{t.Fatalf("same origin rejected: %d",status)}
 p:=a.snapshot()
 if p.Active2!=1||p.Active5!=1||p.Active15!=1||p.Active60!=1||p.Anonymous!=1||p.Foreground!=1||p.LastInteraction==""{t.Fatalf("counts: %+v",p)}
 done:=a.working()
 if p=a.snapshot();p.InFlight!=1{t.Fatalf("work not tracked %+v",p)}
 done()
 if p=a.snapshot();p.InFlight!=0{t.Fatalf("work not released %+v",p)}
 now=now.Add(3*time.Minute)
 if p=a.snapshot();p.Active2!=0||p.Active5!=1{t.Fatalf("rolling expiry wrong %+v",p)}
}

func TestActivityReportIsServerOnlyAggregate(t *testing.T){
 key:=strings.Repeat("k",40)
 var received activityPayload
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Header.Get("Authorization")!="Bearer "+key || r.Header.Get("X-Activity-Service")!="newsgo" {t.Error("service authorization missing")}
  if len(r.Header.Get("X-Activity-Nonce"))<16 || r.Header.Get("X-Activity-Timestamp")==""{t.Error("nonce/time missing")}
  if err:=json.NewDecoder(r.Body).Decode(&received);err!=nil{t.Error(err)}
  w.WriteHeader(204)
 }))
 defer server.Close()
 a,err:=NewActivityReporter(server.URL+"/api/admin/activity/ingest",key)
 if err!=nil{t.Fatal(err)}
 a.mu.Lock()
 a.presence[strings.Repeat("b",32)]=activityPresence{seen:time.Now(),visible:true}
 a.mu.Unlock()
 if err:=a.report(context.Background());err!=nil{t.Fatal(err)}
 if received.Active2!=1 || received.Anonymous!=1 || received.Active60!=1 {t.Fatalf("aggregate %+v",received)}
}

func TestRejectedNewsGoHeartbeatDoesNotBecomeAnonymous(t *testing.T) {
 valid:=strings.Repeat("a",64)
 revoked:=false
 unavailable:=false
 profile:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if unavailable {w.WriteHeader(http.StatusServiceUnavailable);return}
  if revoked {w.WriteHeader(http.StatusUnauthorized);return}
  _,_=w.Write([]byte(`{"user":"guest-test","activitySession":"`+strings.Repeat("1",64)+`","activityUser":"`+strings.Repeat("2",64)+`"}`))
 }))
 defer profile.Close()
 a,_:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("k",40))
 a.profileEndpoint=profile.URL
 a.client=profile.Client()
 ping:=func(bearer string)int{
  r:=httptest.NewRequest(http.MethodPost,"https://newsgo.themauris.org/api/activity/ping",strings.NewReader(`{"id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","visible":true}`))
  r.Header.Set("Origin","https://newsgo.themauris.org");r.Header.Set("Content-Type","application/json")
  if bearer!="" {r.Header.Set("Authorization","Bearer "+bearer)}
  w:=httptest.NewRecorder();a.ping(w,r);return w.Code
 }
 if code:=ping(valid);code!=204{t.Fatalf("valid ping: %d",code)}
 if a.snapshot().Authenticated!=1{t.Fatal("valid profile missing")}
 revoked=true
 if code:=ping(valid);code!=401{t.Fatalf("revoked ping: %d",code)}
 if p:=a.snapshot();p.Authenticated!=0||p.Anonymous!=0{t.Fatalf("revoked identity became anonymous: %+v",p)}
 unavailable=true
 if code:=ping(valid);code!=503{t.Fatalf("validation outage: %d",code)}
 if code:=ping("");code!=204{t.Fatalf("anonymous ping: %d",code)}
}

func TestNewsGoHeartbeatRechecksAccountController(t *testing.T){
 script,err:=os.ReadFile("activity.js");if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(script),"homepage-heartbeat-rejected"){t.Fatal("activity reporter must notify account controller")}
 personal,err:=os.ReadFile("personal_script.go");if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(personal),"void checkHomepageSession();"){t.Fatal("account controller must verify revoked handoff")}
}

func TestActivityScriptIsOptional(t *testing.T) {
 state,err:=New([]feeds.Source{{Name:"Test feed",Category:"World",URL:"https://example.org/rss"}},"https://scrapergo.themauris.org",nil)
 if err!=nil{t.Fatal(err)}
 r:=httptest.NewRequest("GET","http://127.0.0.1:9081/",nil)
 w:=httptest.NewRecorder();state.ServeHTTP(w,r)
 if strings.Contains(w.Body.String(),"/activity.js") {t.Fatal("disabled reporter leaks script")}
 a,_:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("x",40))
 state.SetActivityReporter(a)
 w=httptest.NewRecorder();state.ServeHTTP(w,r)
 if !strings.Contains(w.Body.String(),"/activity.js") {t.Fatal("reporter script missing")}
}

func TestIsolatedActivityHarnessSafetyAndSyntax(t *testing.T) {
 script,err:=os.ReadFile("../../scripts/activity-integration-test.sh")
 if err!=nil{t.Fatal(err)}
 source:=string(script)
 for _,required:=range []string{
  "TEST_ROOT", "homepage-activity.", "127.0.0.1:8788", "127.0.0.1:8085",
  "HOMEPAGE_ACTIVITY_KEY_NEWSGO", "NEWSGO_ACTIVITY_URL", "-settings-dir",
  "LISTEN_ADDR=127.0.0.1:8788",
  "trap cleanup EXIT", "exec env", "lsof", "go test -count=1 ./...",
 } {
  if !strings.Contains(source,required){t.Errorf("missing harness isolation invariant %q",required)}
 }
 for _,forbidden:=range []string{"./newsgo.sh","./homego.sh","cloudflared tunnel","killall","pkill","git reset --hard"}{
  if strings.Contains(source,forbidden){t.Errorf("harness references unsafe operation %q",forbidden)}
 }
 cmd:=exec.Command("bash","-n","../../scripts/activity-integration-test.sh")
 if output,err:=cmd.CombinedOutput();err!=nil{
  t.Fatalf("invalid Bash syntax: %v: %s",err,output)
 }
}

func TestVerifiedNewsGoSessionsDeduplicateTabsAndRetainGuestIdentity(t *testing.T) {
 profiles:=map[string]struct{session,userKey,user string}{
  "a":{strings.Repeat("1",64),strings.Repeat("2",64),"richard"},
  "b":{strings.Repeat("1",64),strings.Repeat("2",64),"richard"},
  "c":{strings.Repeat("3",64),strings.Repeat("2",64),"richard"},
  "d":{strings.Repeat("4",64),strings.Repeat("5",64),"guest-a1b2c3"},
 }
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/session/profile" || r.Header.Get("Origin")!="https://newsgo.themauris.org" {
   t.Errorf("invalid profile request: %s",r.URL.Path);w.WriteHeader(400);return
  }
  bearer:=strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer ")
  entry,ok:=profiles[bearer[:1]]
  if !ok {w.WriteHeader(401);return}
  json.NewEncoder(w).Encode(map[string]string{"activitySession":entry.session,"activityUser":entry.userKey,"user":entry.user})
 }))
 defer server.Close()
 a,_:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("x",40))
 a.profileEndpoint=server.URL+"/api/session/profile"
 a.client=server.Client()
 ping:=func(id,bearer string)int{
  r:=httptest.NewRequest("POST","https://newsgo.themauris.org/api/activity/ping",strings.NewReader(`{"id":"`+id+`","visible":true}`))
  r.Header.Set("Origin","https://newsgo.themauris.org")
  r.Header.Set("Content-Type","application/json")
  r.Header.Set("Authorization","Bearer "+bearer)
  w:=httptest.NewRecorder();a.ping(w,r);return w.Code
 }
 for i,ch:=range []string{"a","b","c","d"}{
  id:=strings.Repeat(string(rune('a'+i)),32)
  if status:=ping(id,ch+strings.Repeat("f",63));status!=204 {t.Fatalf("ping %s: %d",ch,status)}
 }
 p:=a.snapshot()
 if p.Active2!=3 || p.Authenticated!=3 || p.UniqueUsers!=2 || p.Anonymous!=0 {
  t.Fatalf("incorrect session/guest identity counts: %+v",p)
 }
 if len(p.Users)!=2 {t.Fatalf("missing per-user summaries: %+v",p.Users)}
}
