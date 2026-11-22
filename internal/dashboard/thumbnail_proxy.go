package dashboard

import (
 "context"
 "errors"
 "io"
 "log"
 "net"
 "net/http"
 "net/url"
 "os"
 "regexp"
 "strconv"
 "sync"
 "strings"
 "time"
)

var thumbnailServiceTokenPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

func validThumbnailArticleURL(raw string) bool {
 if !safeURL(raw) {return false}
 u,err:=url.Parse(raw)
 if err!=nil {return false}
 host:=strings.TrimSuffix(strings.ToLower(u.Hostname()),".")
 if host=="localhost" || !strings.Contains(host,".") {return false}
 if ip:=net.ParseIP(host);ip!=nil && !publicIP(ip) {return false}
 return true
}

// SetThumbnailBackend configures only the server-to-server thumbnail endpoint.
// Public ScraperGo reader links and classification remain unchanged.
func (s *State) SetThumbnailBackend(raw string) error {
 if raw == "" { s.thumbnailBackend = ""; return nil }
 u,err:=url.Parse(raw)
 if err!=nil || u.Scheme!="http" || u.User!=nil || u.RawQuery!="" || u.Fragment!="" ||
   (u.Path!="" && u.Path!="/") || u.RawPath!="" || u.Opaque!="" ||
   (u.Hostname()!="127.0.0.1" && u.Hostname()!="::1") || u.Port()=="" {
  return errors.New("invalid thumbnail backend: explicit loopback HTTP host and port required")
 }
 s.thumbnailBackend=strings.TrimRight(raw,"/")
 return nil
}

// Classification uses exactly the same strict loopback validation.
func (s *State) SetClassificationBackend(raw string) error {
 var validated State
 if err:=validated.SetThumbnailBackend(raw);err!=nil{return err}
 s.classificationBackend=validated.thumbnailBackend
 return nil
}

const maxThumbnailBytes int64 = 3 << 20

// A shared admission gate honors newly saved limits without replacing
// semaphores while requests are in flight.
var thumbnailProxyGate = struct {
 sync.Mutex
 active int
 peak int
 changed chan struct{}
}{changed:make(chan struct{})}

func thumbnailAcquire(ctx context.Context, wait time.Duration) (time.Duration,error) {
 start:=time.Now()
 timer:=time.NewTimer(wait)
 defer timer.Stop()
 for {
  thumbnailProxyGate.Lock()
  if thumbnailProxyGate.active<currentThumbnailSettings().Concurrency {
   thumbnailProxyGate.active++
   if thumbnailProxyGate.active>thumbnailProxyGate.peak {thumbnailProxyGate.peak=thumbnailProxyGate.active}
   thumbnailProxyGate.Unlock()
   return time.Since(start),nil
  }
  changed:=thumbnailProxyGate.changed
  thumbnailProxyGate.Unlock()
  select {
  case <-changed:
  case <-ctx.Done():return time.Since(start),ctx.Err()
  case <-timer.C:return time.Since(start),context.DeadlineExceeded
  }
 }
}
func thumbnailRelease() {
 thumbnailProxyGate.Lock()
 thumbnailProxyGate.active--
 close(thumbnailProxyGate.changed)
 thumbnailProxyGate.changed=make(chan struct{})
 thumbnailProxyGate.Unlock()
}
func thumbnailActivePeak() (int,int) {
 thumbnailProxyGate.Lock()
 defer thumbnailProxyGate.Unlock()
 return thumbnailProxyGate.active,thumbnailProxyGate.peak
}

// serveThumbnail is a verified-NewsGo-user endpoint (protected by the outer
// Homepage auth gate). It never exposes a ScraperGo service credential to JS.
// ScraperGo grants this shared NewsGo machine credential only scoped API routes.
func (s *State) serveThumbnail(w http.ResponseWriter, r *http.Request) {
 // Capture handler-generated failures separately from the outer auth boundary.
 // Both events share NewsGo's opaque request ID. No article URL or token is logged.
 observed:=&thumbnailResponseCapture{ResponseWriter:w}
 w=observed
 entered:=time.Now()
 defer func(){
  status:=observed.status
  if status==http.StatusOK {return}
  reason:=observed.Header().Get("X-Thumbnail-Outcome")
  if reason=="" {
   if r.Context().Err()!=nil {reason="client_cancelled"} else {reason="handler_rejected"}
  }
  id:=thumbnailRequestID(r)
  if id==""{id="-"}
  log.Printf("THUMBNAIL handler id=%s status=%d outcome=%s duration_ms=%d queue_ms=%s upstream_ms=%s",
   id,status,reason,time.Since(entered).Milliseconds(),
   numericThumbnailMetric(observed.Header().Get("X-Thumbnail-Queue-Ms")),
   numericThumbnailMetric(observed.Header().Get("X-Thumbnail-Upstream-Ms")))
 }()
 // Mark every response produced by this thumbnail handler, including errors.
 // Absence of this marker in Safari is intentionally *unattributed*: an edge,
 // auth gate or other intermediary may have generated or altered the response.
 // No request URLs, identities or credentials are included in these headers.
 w.Header().Set("X-NewsGo-Thumbnail-Proxy","v1")
 w.Header().Set("Cache-Control", "no-store")
 w.Header().Set("X-Content-Type-Options", "nosniff")
 if r.Method != http.MethodGet {
  w.Header().Set("Allow", "GET")
  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  return
 }
 if len(r.URL.RawQuery)>2500 || len(r.URL.Query()["url"])!=1 || len(r.URL.Query()["width"])>1 {
  http.Error(w, "Invalid thumbnail request", http.StatusBadRequest)
  return
 }
 article:=strings.TrimSpace(r.URL.Query().Get("url"))
 if len(article)==0 || len(article)>2048 || !validThumbnailArticleURL(article) {
  http.Error(w, "Invalid article URL", http.StatusBadRequest)
  return
 }
 if width:=r.URL.Query().Get("width"); width!="" && width!="256" {
  http.Error(w, "Unsupported image width", http.StatusBadRequest)
  return
 }
 token:=os.Getenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN")
 if !thumbnailServiceTokenPattern.MatchString(token) {
  // An unconfigured thumbnail backend is not a Homepage verification outage.
  http.Error(w, "Thumbnail service unavailable", http.StatusBadGateway)
  return
 }
 backend:=s.reader
 if s.thumbnailBackend!="" {backend=s.thumbnailBackend}
 reader,err:=url.Parse(backend)
 if err!=nil || reader.User!=nil || reader.RawQuery!="" || reader.Fragment!="" ||
   !(reader.Scheme=="https" || (reader.Scheme=="http" && (reader.Hostname()=="localhost" || reader.Hostname()=="127.0.0.1" || reader.Hostname()=="::1"))) {
  // An unconfigured thumbnail backend is not a Homepage verification outage.
  http.Error(w, "Thumbnail service unavailable", http.StatusBadGateway)
  return
 }
 endpoint:=strings.TrimRight(backend,"/")+"/api/thumbnail?"+url.Values{"url":{article},"width":{"256"}}.Encode()
 // The admission deadline is independent of the upstream fetch budget.
 limits:=currentThumbnailSettings()
 queueDuration,admitErr:=thumbnailAcquire(r.Context(),time.Duration(limits.QueueWaitSeconds)*time.Second)
 queueMS:=strconv.FormatInt(queueDuration.Milliseconds(),10)
 w.Header().Set("X-Thumbnail-Queue-Ms",queueMS)
 w.Header().Set("X-Thumbnail-Active-Limit",strconv.Itoa(limits.Concurrency))
 if admitErr!=nil {
  if errors.Is(admitErr,context.Canceled) || r.Context().Err()!=nil {return}
  w.Header().Set("X-Thumbnail-Outcome","queue_timeout")
  http.Error(w,"Thumbnail queue timeout",http.StatusGatewayTimeout)
  return
 }
 defer thumbnailRelease()
 ctx,cancel:=context.WithTimeout(r.Context(),time.Duration(limits.UpstreamTimeoutSeconds)*time.Second)
 defer cancel()
 upstreamStarted:=time.Now()
 fail:=func(code int,reason string) {
  w.Header().Set("X-Thumbnail-Upstream-Ms",strconv.FormatInt(time.Since(upstreamStarted).Milliseconds(),10))
  w.Header().Set("X-Thumbnail-Outcome",reason)
  http.Error(w,"Thumbnail unavailable",code)
 }
 request,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil)
 if err!=nil {fail(http.StatusBadGateway,"request_error");return}
 request.Header.Set("Authorization","Service "+token)
 request.Header.Set("Accept","image/avif,image/webp,image/png,image/jpeg,image/gif")
 client:=&http.Client{Timeout:time.Duration(limits.UpstreamTimeoutSeconds)*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 response,err:=client.Do(request)
 if err!=nil {
  if errors.Is(err,context.Canceled) && r.Context().Err()!=nil {return}
  if errors.Is(err,context.DeadlineExceeded) || errors.Is(ctx.Err(),context.DeadlineExceeded) {fail(http.StatusGatewayTimeout,"upstream_timeout")}else{fail(http.StatusBadGateway,"transport_error")}
  return
 }
 defer response.Body.Close()
 if response.StatusCode!=http.StatusOK {
  fail(http.StatusBadGateway,"upstream_http_"+strconv.Itoa(response.StatusCode))
  return
 }
 contentType:=strings.ToLower(strings.TrimSpace(strings.Split(response.Header.Get("Content-Type"),";")[0]))
 switch contentType {
 case "image/avif","image/webp","image/png","image/jpeg","image/gif":
 default:
  fail(http.StatusBadGateway,"unsupported_type")
  return
 }
 if response.ContentLength>maxThumbnailBytes {
  fail(http.StatusBadGateway,"oversize")
  return
 }
 data,err:=io.ReadAll(io.LimitReader(response.Body,maxThumbnailBytes+1))
 if err!=nil {
  if errors.Is(err,context.DeadlineExceeded)||errors.Is(ctx.Err(),context.DeadlineExceeded){fail(http.StatusGatewayTimeout,"upstream_timeout")}else{fail(http.StatusBadGateway,"read_error")}
  return
 }
 if int64(len(data))>maxThumbnailBytes {fail(http.StatusBadGateway,"oversize");return}
 sniff:=http.DetectContentType(data)
 if sniff!="image/jpeg" && sniff!="image/png" && sniff!="image/gif" && sniff!="image/webp" && contentType!="image/avif" {
  fail(http.StatusBadGateway,"invalid_image")
  return
 }
 w.Header().Set("X-Thumbnail-Upstream-Ms",strconv.FormatInt(time.Since(upstreamStarted).Milliseconds(),10))
 w.Header().Set("X-Thumbnail-Outcome","loaded")
 w.Header().Set("Content-Type",contentType)
 w.Header().Set("Content-Length",strconv.Itoa(len(data)))
 _,_=w.Write(data)
}

