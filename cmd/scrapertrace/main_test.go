package main

import (
 "bytes"
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestValidateBatchDeduplicatesAndDiscardsFragments(t *testing.T) {
 data:=[]byte(`{"schema":"newsgo.scrapergo-batch.v1","capturedAt":"2026-11-06T00:00:00Z","urls":["https://example.org/story#reader","https://example.org/story","https://other.example/news"]}`)
 urls,err:=validateBatch(data)
 if err!=nil{t.Fatal(err)}
 if len(urls)!=2||urls[0]!="https://example.org/story"{t.Fatalf("unexpected URLs: %#v",urls)}
}

func TestValidateBatchRejectsUnsafeAndMalformedInput(t *testing.T) {
 cases:=[]string{
  `{"schema":"wrong","urls":["https://example.org/story"]}`,
  `{"schema":"newsgo.scrapergo-batch.v1","urls":[]}`,
  `{"schema":"newsgo.scrapergo-batch.v1","urls":["file:///etc/passwd"]}`,
  `{"schema":"newsgo.scrapergo-batch.v1","urls":["https://user:secret@example.org/story"]}`,
  `{"schema":"newsgo.scrapergo-batch.v1","urls":["javascript:alert(1)"]}`,
  `{"schema":"newsgo.scrapergo-batch.v1","urls":["https://example.org/story"],"other":true}`,
 }
 for _,input:=range cases {if _,err:=validateBatch([]byte(input));err==nil{t.Errorf("accepted malformed batch %q",input)}}
}

func TestValidateEndpointLoopbackOnly(t *testing.T) {
 for _,v:=range []string{"https://127.0.0.1:8083/api/read","http://example.com/api/read","http://192.168.1.20:8083/api/read","http://localhost:8083/api/read","http://127.0.0.1:8083/admin","http://127.0.0.1:8083/api/read?token=secret"}{
  if _,err:=validateEndpoint(v);err==nil{t.Errorf("accepted endpoint %s",v)}
 }
 if _,err:=validateEndpoint("http://127.0.0.1:8083/api/read");err!=nil{t.Fatal(err)}
}

func TestReplayUsesNormalReadPOSTAndReportsOnlyIndex(t *testing.T) {
 requested:=[]string{}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!=http.MethodPost||r.URL.Path!="/api/read" {t.Errorf("unexpected request %s %s",r.Method,r.URL.Path)}
  var input map[string]string
  if err:=json.NewDecoder(r.Body).Decode(&input);err!=nil{t.Error(err)}
  requested=append(requested,input["url"])
  if strings.Contains(input["url"],"fail") {http.Error(w,"publisher unavailable",422)}
 }))
 defer server.Close()
 endpoint,err:=validateEndpoint(server.URL+"/api/read")
 if err!=nil{t.Fatal(err)}
 var output bytes.Buffer
 urls:=[]string{"https://example.org/story?secret=one","https://example.org/fail?token=two"}
 success,failed:=replay(context.Background(),urls,endpoint,server.Client(),0,&output)
 if success!=1||failed!=1||len(requested)!=2||requested[0]!=urls[0]||requested[1]!=urls[1]{
  t.Fatalf("unexpected results: success=%d failed=%d received=%#v",success,failed,requested)
 }
 if strings.Contains(output.String(),"secret")||strings.Contains(output.String(),"token")||strings.Contains(output.String(),"story"){
  t.Fatalf("batch output leaked URL: %s",output.String())
 }
}
