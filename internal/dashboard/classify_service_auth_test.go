package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
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
