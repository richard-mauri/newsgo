package dashboard

import (
 "net/http"
 "net/http/httptest"
 "net/url"
 "strings"
 "testing"
)

const testThumbnailToken = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
const testThumbnailArticle = "https://example.com/news/one"

func TestThumbnailProxyRequiresSeparateCredentialAndValidRequest(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN","")
 state:=&State{reader:"https://scrapergo.themauris.org"}
 for _,tc:=range []struct{method,article,width string; want int}{
  {"GET",testThumbnailArticle,"256",http.StatusServiceUnavailable},
  {"POST",testThumbnailArticle,"256",http.StatusMethodNotAllowed},
  {"GET","http://127.0.0.1/private","256",http.StatusBadRequest},
  {"GET","ftp://example.com/a","256",http.StatusBadRequest},
  {"GET",testThumbnailArticle,"999",http.StatusBadRequest},
 }{
  path:="/api/thumbnail?url="+url.QueryEscape(tc.article)+"&width="+tc.width
  w:=httptest.NewRecorder()
  state.serveThumbnail(w,httptest.NewRequest(tc.method,path,nil))
  if w.Code!=tc.want{t.Errorf("%s %q width=%s: HTTP %d, want %d",tc.method,tc.article,tc.width,w.Code,tc.want)}
 }
}

func TestThumbnailProxyPreservesImageAndMachineCredential(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN",testThumbnailToken)
 payload:=[]byte{137,80,78,71,13,10,26,10,0,0,0,0,0,0,0,0}
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path!="/api/thumbnail"{t.Errorf("unexpected path %q",r.URL.Path)}
  if r.Header.Get("Authorization")!="Service "+testThumbnailToken{t.Error("missing thumbnail service authorization")}
  if r.URL.Query().Get("url")!=testThumbnailArticle||r.URL.Query().Get("width")!="256"{t.Error("query mismatch")}
  w.Header().Set("Content-Type","image/png")
  _,_=w.Write(payload)
 }))
 defer upstream.Close()
 state:=&State{reader:upstream.URL}
 w:=httptest.NewRecorder()
 state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle)+"&width=256",nil))
 if w.Code!=http.StatusOK{t.Fatalf("HTTP %d: %s",w.Code,w.Body.String())}
 if w.Header().Get("Content-Type")!="image/png"{t.Errorf("unexpected MIME: %q",w.Header().Get("Content-Type"))}
 if w.Header().Get("Cache-Control")!="no-store"{t.Error("image response must not be stored")}
 if w.Body.String()!=string(payload){t.Error("image bytes changed")}
 if strings.Contains(w.Body.String(),testThumbnailToken){t.Error("service credential exposed")}
}

func TestThumbnailProxyRejectsNonimageAndRedirect(t *testing.T) {
 t.Setenv("SUITE_NEWSGO_SCRAPERGO_THUMBNAIL_TOKEN",testThumbnailToken)
 for _,tc:=range []struct{status int;mime string}{
  {http.StatusFound,"image/png"},
  {http.StatusOK,"text/html"},
  {http.StatusUnauthorized,"text/plain"},
 }{
  upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
   w.Header().Set("Content-Type",tc.mime)
   if tc.status==http.StatusFound{w.Header().Set("Location","https://example.org/should-not-follow")}
   w.WriteHeader(tc.status)
   _,_=w.Write([]byte("not image"))
  }))
  state:=&State{reader:upstream.URL}
  w:=httptest.NewRecorder()
  state.serveThumbnail(w,httptest.NewRequest("GET","/api/thumbnail?url="+url.QueryEscape(testThumbnailArticle),nil))
  if w.Code!=http.StatusBadGateway{t.Errorf("upstream %d %s yielded HTTP %d",tc.status,tc.mime,w.Code)}
  upstream.Close()
 }
}

func TestThumbnailURLUsesProtectedNewsGoRoute(t *testing.T) {
 got:=scraperThumbnailURL("https://scrapergo.themauris.org",testThumbnailArticle)
 if !strings.HasPrefix(got,"/api/thumbnail?"){t.Fatalf("thumbnail URL is not local: %q",got)}
 parsed,err:=url.Parse(got)
 if err!=nil||parsed.Query().Get("url")!=testThumbnailArticle||parsed.Query().Get("width")!="256"{t.Fatalf("bad URL: %q",got)}
 if scraperThumbnailURL("https://scrapergo.themauris.org","file:///etc/passwd")!=""{t.Error("unsafe URL accepted")}
}
