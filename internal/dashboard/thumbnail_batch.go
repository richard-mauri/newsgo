package dashboard

import (
 "context"
 "encoding/base64"
 "encoding/json"
 "io"
 "net/http"
 "net/http/httptest"
 "net/url"
 "log"
 "strconv"
 "time"
 "sync"
)

// Thumbnail batches avoid hundreds of simultaneous browser HTTP admissions.
// Individual images still pass through serveThumbnail's existing security,
// concurrency, upstream and media checks. No new credential path is introduced.
const maxThumbnailBatchSize = 20

type thumbnailBatchRequest struct {
 URLs []string `json:"urls"`
}
type thumbnailBatchItem struct {
 Index int `json:"index"`
 Status int `json:"status"`
 Outcome string `json:"outcome"`
 QueueMS string `json:"queueMs,omitempty"`
 UpstreamMS string `json:"upstreamMs,omitempty"`
 Type string `json:"type,omitempty"`
 Data string `json:"data,omitempty"`
}

func (s *State) serveThumbnailBatch(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Cache-Control","no-store")
 w.Header().Set("X-Content-Type-Options","nosniff")
 if r.Method!=http.MethodPost {w.Header().Set("Allow","POST");http.Error(w,"Method not allowed",http.StatusMethodNotAllowed);return}
 if r.ContentLength>maxThumbnailBatchSize*2300 {http.Error(w,"Thumbnail batch too large",http.StatusRequestEntityTooLarge);return}
 raw,err:=io.ReadAll(io.LimitReader(r.Body,maxThumbnailBatchSize*2300+1))
 if err!=nil||len(raw)>maxThumbnailBatchSize*2300 {http.Error(w,"Thumbnail batch too large",http.StatusRequestEntityTooLarge);return}
 var payload thumbnailBatchRequest
 if json.Unmarshal(raw,&payload)!=nil||len(payload.URLs)==0||len(payload.URLs)>maxThumbnailBatchSize {
  http.Error(w,"Invalid thumbnail batch",http.StatusBadRequest);return
 }
 for _,article:=range payload.URLs {
  if len(article)==0||len(article)>2048||!validThumbnailArticleURL(article){
   http.Error(w,"Invalid thumbnail batch URL",http.StatusBadRequest);return
  }
 }
 // A fixed, small worker pool. The same process-wide admission gate still
 // applies across simultaneous users, including single-image requests.
 concurrency:=currentThumbnailSettings().Concurrency
 if concurrency>6 {concurrency=6}
 if concurrency>len(payload.URLs){concurrency=len(payload.URLs)}
 if concurrency<1 {concurrency=1}
 batchID:=strconv.FormatUint(thumbnailRequestSerial.Add(1),10)
 w.Header().Set("X-NewsGo-Request-ID",batchID)
 w.Header().Set("Content-Type","application/x-ndjson; charset=utf-8")
 w.Header().Set("X-Thumbnail-Batch","v1")
 w.WriteHeader(http.StatusOK)
 if f,ok:=w.(http.Flusher);ok{f.Flush()}
 type job struct{index int;article string}
 jobs:=make(chan job)
 results:=make(chan thumbnailBatchItem,concurrency)
 var wg sync.WaitGroup
 for i:=0;i<concurrency;i++ {
  wg.Add(1)
  go func(){
   defer wg.Done()
   for item:=range jobs {
    if r.Context().Err()!=nil{return}
    // The synthetic request never carries a browser identity or token; it is
    // called only after the enclosing batch has passed Homepage's auth guard.
    sub:=httptest.NewRequest(http.MethodGet,"/api/thumbnail?url="+url.QueryEscape(item.article)+"&width=256",nil).WithContext(context.WithValue(r.Context(),thumbnailTraceIDKey{},batchID+"."+strconv.Itoa(item.index)))
    captured:=httptest.NewRecorder()
    s.serveThumbnail(captured,sub)
    response:=captured.Result()
    result:=thumbnailBatchItem{Index:item.index,Status:response.StatusCode,Outcome:response.Header.Get("X-Thumbnail-Outcome"),
     QueueMS:response.Header.Get("X-Thumbnail-Queue-Ms"),UpstreamMS:response.Header.Get("X-Thumbnail-Upstream-Ms")}
    if response.StatusCode==http.StatusOK {
     result.Type=response.Header.Get("Content-Type")
     result.Data=base64.StdEncoding.EncodeToString(captured.Body.Bytes())
    }
    select{case results<-result:case <-r.Context().Done():return}
   }
  }()
 }
 go func(){
  defer close(jobs)
  for i,article:=range payload.URLs {
   select{case jobs<-job{i,article}:case <-r.Context().Done():return}
  }
 }()
 go func(){wg.Wait();close(results)}()
 encoder:=json.NewEncoder(w)
 started:=time.Now()
 completed,failed:=0,0
 for item:=range results {
  if r.Context().Err()!=nil{return}
  if err:=encoder.Encode(item);err!=nil{return}
  completed++
  if item.Status!=http.StatusOK{failed++}
  if f,ok:=w.(http.Flusher);ok{f.Flush()}
 }
 log.Printf("THUMBNAIL batch id=%s requested=%d completed=%d failed=%d duration_ms=%d",batchID,len(payload.URLs),completed,failed,time.Since(started).Milliseconds())
}
