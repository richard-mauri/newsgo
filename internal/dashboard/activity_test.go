package dashboard

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
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

func TestActivityScriptIsOptional(t *testing.T) {
 state,err:=New(nil,"https://scrapergo.themauris.org",nil)
 if err!=nil{t.Fatal(err)}
 r:=httptest.NewRequest("GET","http://127.0.0.1:9081/",nil)
 w:=httptest.NewRecorder();state.ServeHTTP(w,r)
 if strings.Contains(w.Body.String(),"/activity.js") {t.Fatal("disabled reporter leaks script")}
 a,_:=NewActivityReporter("http://127.0.0.1:8788/api/admin/activity/ingest",strings.Repeat("x",40))
 state.SetActivityReporter(a)
 w=httptest.NewRecorder();state.ServeHTTP(w,r)
 if !strings.Contains(w.Body.String(),"/activity.js") {t.Fatal("reporter script missing")}
}
