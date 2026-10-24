package dashboard

import (
 "bytes"
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestAdminClassificationDiagnosticRequiresAuthentication(t *testing.T) {
 state:=&State{}
 req:=httptest.NewRequest(http.MethodPost,"/api/admin/classification-diagnostic",strings.NewReader(`{"url":"https://example.org/article"}`))
 rec:=httptest.NewRecorder()
 state.serveAdminClassificationDiagnostic(rec,req)
 if rec.Code!=http.StatusServiceUnavailable {t.Fatalf("want 503 without sync store, got %d",rec.Code)}
}

func TestFreshClassificationDiagnostics(t *testing.T) {
 for _,tt:=range []struct{name string;code int;body string;want string;classification string}{
  {"article",200,`{"classification":"article"}`,"","article"},
  {"upstream-error",502,"unavailable","http_error","unknown"},
  {"bad-json",200,"invalid","invalid_json","unknown"},
  {"unexpected",200,`{"classification":"unexpected"}`,"unexpected_classification","unknown"},
 } {
  t.Run(tt.name,func(t *testing.T){
   remote:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
    if r.Method!="POST"||r.URL.Path!="/api/classify" {t.Errorf("unexpected request: %s %s",r.Method,r.URL.Path)}
    w.WriteHeader(tt.code)
    _,_=w.Write([]byte(tt.body))
   }))
   defer remote.Close()
   got:=classifyDiagnostic(context.Background(),remote.URL,"https://example.org/article")
   if got.Failure!=tt.want ||got.Classification!=tt.classification||got.HTTPStatus!=tt.code {
    t.Fatalf("unexpected diagnostic: %+v",got)
   }
  })
 }
}

func TestAdminClassificationRejectsNonPOST(t *testing.T) {
 state:=&State{}
 req:=httptest.NewRequest(http.MethodGet,"/api/admin/classification-diagnostic",bytes.NewReader(nil))
 rec:=httptest.NewRecorder()
 state.serveAdminClassificationDiagnostic(rec,req)
 if rec.Code!=http.StatusMethodNotAllowed {t.Fatalf("got %d",rec.Code)}
}
