package dashboard

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestClassifyDiagnosticPublisherHTTPStatus(t *testing.T) {
 cases:=[]struct{name,body,classification,reason string;want int}{
  {"forbidden",`{"classification":"unavailable","reason":"publisher_unavailable","publisherHttpStatus":403}`,"unavailable","publisher_unavailable",403},
  {"unauthorized",`{"classification":"unavailable","reason":"publisher_unavailable","publisherHttpStatus":401}`,"unavailable","publisher_unavailable",401},
  {"rate limit",`{"classification":"unavailable","reason":"publisher_unavailable","publisherHttpStatus":429}`,"unavailable","publisher_unavailable",429},
  {"successful",`{"classification":"article"}`,"article","",0},
  {"legacy",`{"classification":"unavailable","reason":"publisher_unavailable"}`,"unavailable","publisher_unavailable",0},
 }
 for _,tc:=range cases{t.Run(tc.name,func(t *testing.T){
  server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
   if r.URL.Path!="/api/classify"{t.Errorf("unexpected path %q",r.URL.Path)}
   w.Header().Set("Content-Type","application/json")
   _,_=w.Write([]byte(tc.body))
  }))
  defer server.Close()
  result:=classifyDiagnostic(context.Background(),server.URL,"https://example.com/article")
  if result.Classification!=tc.classification||result.PublisherHTTPStatus!=tc.want||result.Reason!=tc.reason||result.Failure!=""{
   t.Fatalf("result %+v",result)
  }
 })}
}

func TestAdminDiagnosticScriptPublisherStatus(t *testing.T) {
 required:=[]string{"Publisher HTTP status: ","publisherHttpStatus","lines.join('\\n')"}
 for _,item:=range required {if !strings.Contains(personalJS,item){t.Errorf("missing %q",item)}}
 for _,item:=range []string{"'\\\\nHTTP status: '","'\\\\nTimeout: '"} {
  if strings.Contains(personalJS,item){t.Errorf("misleading old field %q",item)}
 }
}

