package dashboard

import (
 "bufio"
 "bytes"
 "encoding/base64"
 "encoding/json"
 "fmt"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync/atomic"
 "testing"
 "time"
)

func TestThumbnailBatchValidatesBoundsAndMethods(t *testing.T){
 s:=&State{reader:"https://scrapergo.themauris.org"}
 for _,tc:=range []struct{name,method,body string;want int}{
  {"get","GET","",405},
  {"empty","POST",`{"urls":[]}`,400},
  {"malformed","POST","[",400},
  {"private-host","POST",`{"urls":["http://127.0.0.1/private"]}`,400},
  {"overlong","POST",fmt.Sprintf(`{"urls":[%s]}`,strings.TrimSuffix(strings.Repeat(`"https://example.com",`,21),",")),400},
  {"invalid-type","POST",`{"urls":[123]}`,400},
 }{
  t.Run(tc.name,func(t *testing.T){
   r:=httptest.NewRequest(tc.method,"/api/thumbnail/batch",strings.NewReader(tc.body))
   w:=httptest.NewRecorder()
   s.serveThumbnailBatch(w,r)
   if w.Code!=tc.want{t.Fatalf("got %d want %d",w.Code,tc.want)}
  })
 }
}

func TestThumbnailBatchStreamsMixedOutOfOrderResults(t *testing.T){
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN",testThumbnailToken)
 var active atomic.Int32
 var peak atomic.Int32
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  cur:=active.Add(1)
  for{old:=peak.Load();if cur<=old||peak.CompareAndSwap(old,cur){break}}
  defer active.Add(-1)
  article:=r.URL.Query().Get("url")
  if strings.Contains(article,"slow"){time.Sleep(250*time.Millisecond)}
  if strings.Contains(article,"missing"){w.WriteHeader(404);return}
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write([]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0})
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 urls:=[]string{"https://example.com/slow","https://example.com/missing","https://example.com/fast"}
 payload,_:=json.Marshal(thumbnailBatchRequest{URLs:urls})
 w:=httptest.NewRecorder()
 state.serveThumbnailBatch(w,httptest.NewRequest("POST","/api/thumbnail/batch",bytes.NewReader(payload)))
 if w.Code!=200||w.Header().Get("X-Thumbnail-Batch")!="v1" {t.Fatalf("bad batch headers %d",w.Code)}
 scan:=bufio.NewScanner(bytes.NewReader(w.Body.Bytes()))
 scan.Buffer(make([]byte,1024),int(maxThumbnailBytes*2))
 var items []thumbnailBatchItem
 for scan.Scan(){
  var item thumbnailBatchItem
  if e:=json.Unmarshal(scan.Bytes(),&item);e!=nil{t.Fatal(e)}
  items=append(items,item)
 }
 if e:=scan.Err();e!=nil{t.Fatal(e)}
 if len(items)!=3{t.Fatalf("got %d items",len(items))}
 if items[0].Index==0{t.Fatal("slow item blocked fast streamed result")}
 statuses:=map[int]int{}
 for _,item:=range items{
  statuses[item.Index]=item.Status
  if item.Index==1{
   if item.Status!=502||item.Outcome!="upstream_http_404"||item.Data!=""{t.Fatalf("bad 404 item: %+v",item)}
  }else{
   data,e:=base64.StdEncoding.DecodeString(item.Data)
   if e!=nil||len(data)!=16||item.Status!=200||item.Type!="image/png"{t.Fatalf("bad image item %d: %v",item.Index,e)}
  }
 }
 if statuses[0]!=200||statuses[1]!=502||statuses[2]!=200{t.Fatalf("bad statuses %v",statuses)}
 if peak.Load()>6 {t.Fatalf("unbounded worker concurrency %d",peak.Load())}
}

func TestThumbnailBatchRemainsProtectedByHomepageGuard(t *testing.T){
 app:=NewNewsGoAuthGate().Protect(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/api/thumbnail/batch"{t.Fatal("unverified request reached batch handler")}
 }))
 payload:=strings.NewReader(`{"urls":["https://example.com/a"]}`)
 w:=httptest.NewRecorder()
 app.ServeHTTP(w,httptest.NewRequest("POST","/api/thumbnail/batch",payload))
 if w.Code!=http.StatusUnauthorized{t.Fatalf("unauthorized batch status %d",w.Code)}
}
