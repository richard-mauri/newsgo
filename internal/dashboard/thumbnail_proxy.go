package dashboard

import (
 "context"
 "errors"
 "io"
 "net"
 "net/http"
 "net/url"
 "os"
 "regexp"
 "strconv"
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

const maxThumbnailBytes int64 = 3 << 20

var thumbnailProxySlots = make(chan struct{}, 6)

// serveThumbnail is a verified-NewsGo-user endpoint (protected by the outer
// Homepage auth gate). It never exposes a ScraperGo service credential to JS.
// ScraperGo must separately grant GET /api/thumbnail to this distinct credential.
func (s *State) serveThumbnail(w http.ResponseWriter, r *http.Request) {
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
 token:=os.Getenv("SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN")
 if !thumbnailServiceTokenPattern.MatchString(token) || token==os.Getenv("SUITE_NEWSGO_SCRAPERGO_CLASSIFY_TOKEN") {
  http.Error(w, "Thumbnail service unavailable", http.StatusServiceUnavailable)
  return
 }
 reader,err:=url.Parse(s.reader)
 if err!=nil || reader.User!=nil || reader.RawQuery!="" || reader.Fragment!="" ||
   !(reader.Scheme=="https" || (reader.Scheme=="http" && (reader.Hostname()=="localhost" || reader.Hostname()=="127.0.0.1" || reader.Hostname()=="::1"))) {
  http.Error(w, "Thumbnail service unavailable", http.StatusServiceUnavailable)
  return
 }
 endpoint:=strings.TrimRight(s.reader,"/")+"/api/thumbnail?"+url.Values{"url":{article},"width":{"256"}}.Encode()
 ctx,cancel:=context.WithTimeout(r.Context(),12*time.Second)
 defer cancel()
 request,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil)
 if err!=nil {http.Error(w,"Thumbnail unavailable",http.StatusBadGateway);return}
 request.Header.Set("Authorization","Service "+token)
 request.Header.Set("Accept","image/avif,image/webp,image/png,image/jpeg,image/gif")
 client:=&http.Client{Timeout:12*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 select {
 case thumbnailProxySlots<-struct{}{}:defer func(){<-thumbnailProxySlots}()
 case <-ctx.Done():return
 }
 response,err:=client.Do(request)
 if err!=nil {
  if errors.Is(err,context.Canceled) {return}
  http.Error(w,"Thumbnail unavailable",http.StatusBadGateway)
  return
 }
 defer response.Body.Close()
 if response.StatusCode!=http.StatusOK {
  http.Error(w,"Thumbnail unavailable",http.StatusBadGateway)
  return
 }
 contentType:=strings.ToLower(strings.TrimSpace(strings.Split(response.Header.Get("Content-Type"),";")[0]))
 switch contentType {
 case "image/avif","image/webp","image/png","image/jpeg","image/gif":
 default:
  http.Error(w,"Unsupported image type",http.StatusBadGateway)
  return
 }
 if response.ContentLength>maxThumbnailBytes {
  http.Error(w,"Thumbnail too large",http.StatusBadGateway)
  return
 }
 data,err:=io.ReadAll(io.LimitReader(response.Body,maxThumbnailBytes+1))
 if err!=nil || int64(len(data))>maxThumbnailBytes {
  http.Error(w,"Thumbnail too large",http.StatusBadGateway)
  return
 }
 // Header declarations alone do not establish that the body is image data.
 sniff:=http.DetectContentType(data)
 if sniff!="image/jpeg" && sniff!="image/png" && sniff!="image/gif" && sniff!="image/webp" && contentType!="image/avif" {
  http.Error(w,"Invalid image response",http.StatusBadGateway)
  return
 }
 w.Header().Set("Content-Type",contentType)
 w.Header().Set("Content-Length",strconv.Itoa(len(data)))
 _,_=w.Write(data)
}

