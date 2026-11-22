package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

func TestClassificationServiceHeaderOnPipelineAndDiagnostic(t *testing.T){
 key:=strings.Repeat("b",64)
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",key)
 count:=0
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  count++
  if r.Method!=http.MethodPost||r.URL.Path!="/api/classify"{t.Errorf("bad service route: %s %s",r.Method,r.URL.Path)}
  if r.Header.Get("Authorization")!="Service "+key{t.Error("missing scoped service credential")}
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article","readableArticle":true}`))
 }))
 defer server.Close()
 // Disable cached results for a deterministic outgoing request.
 if got:=newsClassifyWithCache(context.Background(),server.URL,"https://example.org/new-unique-service-check",false);got!="article"{t.Fatalf("pipeline got %q",got)}
 diag:=classifyDiagnostic(context.Background(),server.URL,"https://example.org/new-unique-diagnostic-check")
 if diag.Classification!="article"{t.Fatalf("diagnostic: %+v",diag)}
 if count!=2{t.Fatalf("expected two authenticated server calls, got %d",count)}
}

func TestClassificationCredentialNeverLeavesAllowedTransport(t *testing.T){
 key:=strings.Repeat("c",64)
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",key)
 for _,u:=range []string{"http://outside.example/api/classify","https://outside.example/api/read","https://name:pass@outside.example/api/classify"}{
  r,_:=http.NewRequest("POST",u,nil)
  if err:=authorizeScraperGoClassification(r);err==nil{t.Fatalf("must reject %s",u)}
  if r.Header.Get("Authorization")!=""{t.Fatalf("credential leaked on %s",u)}
 }
 r,_:=http.NewRequest("POST","https://scrapergo.themauris.org/api/classify",nil)
 if err:=authorizeScraperGoClassification(r);err!=nil{t.Fatal(err)}
 if r.Header.Get("Authorization")!="Service "+key{t.Fatal("unexpected service header")}
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN","malformed")
 r2,_:=http.NewRequest("POST","https://scrapergo.themauris.org/api/classify",nil)
 if authorizeScraperGoClassification(r2)==nil{t.Fatal("malformed configuration accepted")}
}

func TestClassificationLoopbackKeepsPublicCacheIdentity(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",strings.Repeat("d",64))
 const publicReader="https://scrapergo.themauris.org"
 const article="https://example.org/loopback-classification"
 count:=0
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  count++
  if r.Method!=http.MethodPost||r.URL.Path!="/api/classify"{t.Errorf("wrong route %s %s",r.Method,r.URL.Path)}
  if r.Header.Get("Authorization")!="Service "+strings.Repeat("d",64){t.Error("missing service credential")}
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer upstream.Close()
 state:=&State{reader:publicReader}
 if err:=state.SetClassificationBackend(upstream.URL);err!=nil{t.Fatal(err)}
 if state.reader!=publicReader {t.Fatal("public reader mutated")}
 if state.classificationBackend==publicReader {t.Fatal("classification not separated")}
 newsClassificationCache.Lock()
 original,hadOriginal:=newsClassificationCache.entries[publicReader+"|"+article]
 delete(newsClassificationCache.entries,publicReader+"|"+article)
 newsClassificationCache.Unlock()
 defer func(){
  newsClassificationCache.Lock()
  if hadOriginal {newsClassificationCache.entries[publicReader+"|"+article]=original
  } else {delete(newsClassificationCache.entries,publicReader+"|"+article)}
  newsClassificationCache.Unlock()
 }()
 if got:=newsClassifyWithCacheBackendStats(context.Background(),state.reader,state.classificationBackend,article,false,nil);got!="article"{t.Fatal(got)}
 // The first request bypasses the cache; populate an in-memory entry under
 // the public reader key as an existing installation would have done.
 newsClassificationCache.Lock()
 newsClassificationCache.entries[publicReader+"|"+article]=newsClassificationEntry{value:"article",expiry:time.Now().Add(time.Minute)}
 newsClassificationCache.Unlock()
 if got:=newsClassifyWithCacheBackendStats(context.Background(),state.reader,state.classificationBackend,article,true,nil);got!="article"{t.Fatal(got)}
 if count!=1 {t.Fatalf("cache identity changed; upstream calls=%d",count)}
 diagnostic:=classifyDiagnostic(context.Background(),classificationRequestBackend(state.reader,state.classificationBackend),article)
 if diagnostic.Classification!="article" || count!=2 {t.Fatalf("uncached diagnostic failed: %+v count %d",diagnostic,count)}
}

func TestClassificationLoopbackRejectsUnsafeOverride(t *testing.T){
 state:=&State{reader:"https://scrapergo.themauris.org"}
 for _,raw:=range []string{"http://localhost:8083","https://127.0.0.1:8083","http://10.0.0.1:8083","http://127.0.0.1:8083/api/read","http://x:y@127.0.0.1:8083"}{
  if err:=state.SetClassificationBackend(raw);err==nil{t.Errorf("accepted unsafe backend %q",raw)}
  if state.classificationBackend!="" {t.Error("unsafe backend retained")}
 }
 if err:=state.SetClassificationBackend("http://127.0.0.1:8083");err!=nil{t.Fatal(err)}
 if err:=state.SetClassificationBackend("");err!=nil||state.classificationBackend!=""{t.Fatal("override not cleared")}
}
