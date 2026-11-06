// scrapertrace replays a NewsGo-exported batch through a locally running
// ScraperGo service, using its normal URL validation and access policy.
package main

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "io"
 "net"
 "net/http"
 "net/url"
 "os"
 "strings"
 "time"
)

const batchSchema="newsgo.scrapergo-batch.v1"
const maxBatchURLs=500

type batchFile struct {
 Schema string `json:"schema"`
 CapturedAt string `json:"capturedAt"`
 URLs []string `json:"urls"`
}

func validateBatch(data []byte) ([]string,error) {
 var b batchFile
 if len(data)>1<<20{return nil,errors.New("batch file exceeds 1 MiB")}
 dec:=json.NewDecoder(bytes.NewReader(data))
 dec.DisallowUnknownFields()
 if err:=dec.Decode(&b);err!=nil{return nil,fmt.Errorf("invalid batch JSON: %w",err)}
 var trailing any
 if err:=dec.Decode(&trailing);err!=io.EOF{return nil,errors.New("batch must contain one JSON object")}
 if b.Schema!=batchSchema{return nil,fmt.Errorf("unexpected batch schema")}
 if len(b.URLs)==0||len(b.URLs)>maxBatchURLs{return nil,fmt.Errorf("expected 1–%d URLs",maxBatchURLs)}
 out:=make([]string,0,len(b.URLs));seen:=map[string]bool{}
 for _,raw:=range b.URLs {
  if len(raw)>2048||strings.ContainsAny(raw,"\r\n\t"){return nil,errors.New("invalid article URL")}
  u,err:=url.Parse(raw)
  if err!=nil||u==nil||(u.Scheme!="http"&&u.Scheme!="https")||u.Hostname()==""||u.User!=nil{
   return nil,errors.New("batch contains invalid HTTP(S) URL")
  }
  u.Fragment=""
  normalized:=u.String()
  if seen[normalized]{continue}
  seen[normalized]=true;out=append(out,normalized)
 }
 return out,nil
}

func validateEndpoint(raw string) (string,error) {
 u,err:=url.Parse(raw)
 if err!=nil||u==nil||u.Scheme!="http"||u.User!=nil||u.RawQuery!=""||u.Fragment!=""||u.Path!="/api/read"{
  return "",errors.New("ScraperGo endpoint must be a loopback HTTP /api/read address")
 }
 ip:=net.ParseIP(u.Hostname())
 if ip==nil||!ip.IsLoopback(){return "",errors.New("ScraperGo endpoint must use a numeric loopback address")}
 return u.String(),nil
}

// replay records classification and reader results independently. Classification
// invokes ScraperGo's normal reader internally, but may be cached or time out.
type probeResult struct {
 HTTPStatus int
 Classification string
 Kind string
 PublisherHTTPStatus int
 RequestError bool
}

func endpointForRead(readEndpoint,path string) string {
 u,_:=url.Parse(readEndpoint)
 u.Path=path
 return u.String()
}

func probe(ctx context.Context,address,endpoint string,client *http.Client,isClassification bool) probeResult {
 body,_:=json.Marshal(map[string]string{"url":address})
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,endpoint,bytes.NewReader(body))
 if err!=nil{return probeResult{RequestError:true}}
 req.Header.Set("Content-Type","application/json")
 response,err:=client.Do(req)
 if err!=nil{return probeResult{RequestError:true}}
 defer response.Body.Close()
 result:=probeResult{HTTPStatus:response.StatusCode}
 // Bound response inspection and decode only documented, non-sensitive fields.
 limited,err:=io.ReadAll(io.LimitReader(response.Body,8193))
 if err!=nil||len(limited)>8192{return result}
 var decoded struct {
  Classification string `json:"classification"`
  Kind string `json:"kind"`
  PublisherHTTPStatus int `json:"publisherHttpStatus"`
 }
 if json.Unmarshal(limited,&decoded)!=nil{return result}
 if isClassification {
  switch decoded.Classification {
  case "article","homepage","unreadable","unavailable","unknown":
   result.Classification=decoded.Classification
  }
 } else {
  // Restrict free-form upstream content to a small stable diagnostic vocabulary.
  switch decoded.Kind {
  case "page_index","extraction_error","publisher_access_denied","access_error","pdf_download_failed","pdf_extraction_failed":
   result.Kind=decoded.Kind
  }
  if decoded.PublisherHTTPStatus>=100&&decoded.PublisherHTTPStatus<=599 {
   result.PublisherHTTPStatus=decoded.PublisherHTTPStatus
  }
 }
 return result
}

func replay(ctx context.Context, urls []string, endpoint string, client *http.Client, pace time.Duration, output io.Writer) (int,int) {
 succeeded,failed,mismatched:=0,0,0
 classificationEndpoint:=endpointForRead(endpoint,"/api/classify")
 for i,address:=range urls {
  if i>0&&pace>0{
   select {case <-ctx.Done():return succeeded,failed+len(urls)-i;case <-time.After(pace):}
  }
  classification:=probe(ctx,address,classificationEndpoint,client,true)
  read:=probe(ctx,address,endpoint,client,false)
  classLabel:=classification.Classification
  if classLabel==""{classLabel="unreported"}
  classHTTP:=fmt.Sprintf("%d",classification.HTTPStatus)
  if classification.RequestError{classHTTP="request_failed"}
  readHTTP:=fmt.Sprintf("%d",read.HTTPStatus)
  if read.RequestError{readHTTP="request_failed"}
  disagreement:=classification.Classification=="article"&&(read.RequestError||read.HTTPStatus<200||read.HTTPStatus>=300)
  if disagreement{mismatched++}
  if !read.RequestError&&read.HTTPStatus>=200&&read.HTTPStatus<300{succeeded++}else{failed++}
  fmt.Fprintf(output,"Article %d/%d: classification=%s classify_http=%s read_http=%s",i+1,len(urls),classLabel,classHTTP,readHTTP)
  if read.Kind!=""{fmt.Fprintf(output," read_kind=%s",read.Kind)}
  if read.PublisherHTTPStatus!=0{fmt.Fprintf(output," publisher_http=%d",read.PublisherHTTPStatus)}
  if disagreement{fmt.Fprint(output," mismatch=classified_article_but_unreadable")}
  fmt.Fprintln(output)
 }
 fmt.Fprintf(output,"Classification/read disagreements (article classified, read failed): %d\n",mismatched)
 return succeeded,failed
}

func main() {
 input:=flag.String("input","","path to NewsGo diagnostic URL export")
 pace:=flag.Duration("delay",time.Second,"delay between article requests")
 endpointFlag:=flag.String("scrapergo","http://127.0.0.1:8083/api/read","local ScraperGo read endpoint (loopback only)")
 flag.Parse()
 if *input==""||*pace<0||*pace>time.Minute||flag.NArg()!=0{
  fmt.Fprintln(os.Stderr,"Usage: go run ./cmd/scrapertrace -input ~/Downloads/newsgo-scrapergo-urls.json [-delay 1s]")
  os.Exit(2)
 }
 endpoint,err:=validateEndpoint(*endpointFlag)
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(2)}
 f,err:=os.Open(*input)
 if err!=nil{fmt.Fprintln(os.Stderr,"Cannot open batch:",err);os.Exit(1)}
 data,err:=io.ReadAll(io.LimitReader(f,1<<20+1))
 _=f.Close()
 if err!=nil{fmt.Fprintln(os.Stderr,"Cannot read batch:",err);os.Exit(1)}
 urls,err:=validateBatch(data)
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(2)}
 client:=&http.Client{Timeout:45*time.Second,CheckRedirect:func(_ *http.Request,_ []*http.Request)error{return http.ErrUseLastResponse}}
 fmt.Printf("Classifying and reading %d distinct article URLs with local ScraperGo; pause %s between URL pairs.\n",len(urls),*pace)
 good,bad:=replay(context.Background(),urls,endpoint,client,*pace,os.Stdout)
 fmt.Printf("Finished: %d succeeded, %d failed. Inspect ScraperGo: go run ./cmd/recoveryreport (from its repository).\n",good,bad)
 if bad>0{os.Exit(1)}
}
