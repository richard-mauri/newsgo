package main

import (
 "bytes"
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "sync/atomic"
 "time"
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

func TestReplayReportsClassificationAndReadStatusWithoutURLs(t *testing.T) {
 var classify,read []string
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!=http.MethodPost {t.Errorf("wrong method: %s",r.Method)}
  var input struct{URL string `json:"url"`}
  if err:=json.NewDecoder(r.Body).Decode(&input);err!=nil{t.Error(err)}
  w.Header().Set("Content-Type","application/json")
  switch r.URL.Path {
  case "/api/classify":
   classify=append(classify,input.URL)
   if strings.Contains(input.URL,"bad") {
    _,_=w.Write([]byte(`{"classification":"article","readableArticle":true}`))
   } else {
    _,_=w.Write([]byte(`{"classification":"article","readableArticle":true}`))
   }
  case "/api/read":
   read=append(read,input.URL)
   if strings.Contains(input.URL,"bad"){
    w.WriteHeader(422)
    _,_=w.Write([]byte(`{"kind":"publisher_access_denied","publisherHttpStatus":403,"error":"contains-private-data"}`))
   }else{
    _,_=w.Write([]byte(`{"title":"Test article"}`))
   }
  default:t.Errorf("unexpected endpoint %s",r.URL.Path)
  }
 }))
 defer server.Close()
 endpoint,err:=validateEndpoint(server.URL+"/api/read")
 if err!=nil{t.Fatal(err)}
 urls:=[]string{"https://example.org/story?secret=one","https://example.org/bad?token=two"}
 var output bytes.Buffer
 success,failed:=replay(context.Background(),urls,endpoint,server.Client(),0,&output)
 if success!=1||failed!=1||len(classify)!=2||len(read)!=2{t.Fatalf("unexpected totals: successes=%d failures=%d classify=%d read=%d",success,failed,len(classify),len(read))}
 for i:=range urls {if classify[i]!=urls[i]||read[i]!=urls[i]{t.Fatal("wrong URL submitted")}}
 result:=output.String()
 for _,need:=range []string{"Article 1/2: classification=article classify_http=200 read_http=200","Article 2/2: classification=article classify_http=200 read_http=422 read_kind=publisher_access_denied publisher_http=403 mismatch=classified_article_but_unreadable","Classification/read disagreements (article classified, read failed): 1"} {
  if !strings.Contains(result,need){t.Errorf("missing output %q from %s",need,result)}
 }
 for _,private:=range []string{"secret","token","example.org","contains-private-data"}{
  if strings.Contains(result,private){t.Errorf("leaked private content %q",private)}
 }
}

func TestReplayHandlesUnclassifiedAndNonJSONResponses(t *testing.T){
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api/classify"{
   w.WriteHeader(422)
   _,_=w.Write([]byte("not json"))
   return
  }
  w.WriteHeader(422)
  _,_=w.Write([]byte(`{"kind":"extraction_error","error":"secret"}`))
 }))
 defer server.Close()
 var output bytes.Buffer
 good,bad:=replay(context.Background(),[]string{"https://example.org/private"},server.URL+"/api/read",server.Client(),0,&output)
 if good!=0||bad!=1 {t.Fatalf("unexpected totals %d %d",good,bad)}
 if !strings.Contains(output.String(),"classification=unreported classify_http=422 read_http=422 read_kind=extraction_error"){t.Fatal(output.String())}
 if strings.Contains(output.String(),"secret")||strings.Contains(output.String(),"example.org"){t.Fatal("diagnostic leaked content")}
}

func TestProbeRestrictsUntrustedResponseLabels(t *testing.T) {
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api/classify"{
   _,_=w.Write([]byte(`{"classification":"article\\nINJECTED"}`))
  }else{
   w.WriteHeader(422)
   _,_=w.Write([]byte(`{"kind":"custom\\nINJECTED","publisherHttpStatus":999}`))
  }
 }))
 defer server.Close()
 var out bytes.Buffer
 _,_=replay(context.Background(),[]string{"https://example.org/test"},server.URL+"/api/read",server.Client(),0,&out)
 if strings.Contains(out.String(),"INJECTED")||strings.Contains(out.String(),"publisher_http="){t.Fatal(out.String())}
}

func TestReplayConcurrentTimingAndBound(t *testing.T) {
 var active,maxActive int64
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  current:=atomic.AddInt64(&active,1)
  for {
   previous:=atomic.LoadInt64(&maxActive)
   if current<=previous||atomic.CompareAndSwapInt64(&maxActive,previous,current){break}
  }
  time.Sleep(20*time.Millisecond)
  atomic.AddInt64(&active,-1)
  w.Header().Set("Content-Type","application/json")
  if r.URL.Path=="/api/classify"{
   _,_=w.Write([]byte(`{"classification":"article"}`))
  }else if r.URL.Path=="/api/read"{
   _,_=w.Write([]byte(`{"title":"ok"}`))
  }else{t.Errorf("unexpected path %s",r.URL.Path)}
 }))
 defer server.Close()
 urls:=[]string{
  "https://example.org/1","https://example.org/2",
  "https://example.org/3","https://example.org/4",
  "https://example.org/5","https://example.org/6",
 }
 var output bytes.Buffer
 good,bad:=replayWithWorkers(context.Background(),urls,server.URL+"/api/read",server.Client(),0,4,&output)
 if good!=len(urls)||bad!=0{t.Fatalf("bad totals good=%d bad=%d",good,bad)}
 if atomic.LoadInt64(&maxActive)<2||atomic.LoadInt64(&maxActive)>4{t.Fatalf("unexpected concurrency=%d",maxActive)}
 for _,wanted:=range []string{"classify_wall=","read_wall=","Classification wall latency: n=6","Read wall latency: n=6","Batch elapsed wall time:"}{
  if !strings.Contains(output.String(),wanted){t.Errorf("missing %q in %s",wanted,output.String())}
 }
}

func TestReplayDispatchDelayExcludedFromLatency(t *testing.T){
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api/classify"{_,_=w.Write([]byte(`{"classification":"article"}`));return}
  _,_=w.Write([]byte(`{"title":"ok"}`))
 }))
 defer server.Close()
 urls:=[]string{"https://example.org/1","https://example.org/2"}
 var output bytes.Buffer
 started:=time.Now()
 good,bad:=replayWithWorkers(context.Background(),urls,server.URL+"/api/read",server.Client(),100*time.Millisecond,1,&output)
 if good!=2||bad!=0{t.Fatalf("bad counts: %d %d",good,bad)}
 if time.Since(started)<100*time.Millisecond{t.Fatal("dispatch throttle omitted")}
 if !strings.Contains(output.String(),"dispatch_delay=100ms"){t.Fatal(output.String())}
}

func TestLatencySummaryPercentiles(t *testing.T){
 values:=[]time.Duration{5*time.Millisecond,1*time.Millisecond,4*time.Millisecond,2*time.Millisecond,3*time.Millisecond}
 result:=latencySummary(values)
 for _,item:=range []string{"n=5","min=1ms","mean=3ms","p50=3ms","p95=5ms","max=5ms"}{
  if !strings.Contains(result,item){t.Errorf("missing %s in %s",item,result)}
 }
}
