package dashboard

import (
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 suiteui "github.com/richard-mauri/homepage/shared/ui"
)

var suiteButtonCSS=string(suiteui.ButtonCSS)

func TestCanonicalButtonStylesheetServedLocally(t *testing.T) {
 s:=&State{}
 r:=httptest.NewRequest(http.MethodGet,"/suite-ui/buttons.css",nil)
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,r)
 if w.Code!=http.StatusOK {t.Fatalf("CSS status %d",w.Code)}
 if got:=w.Header().Get("Content-Type");got!="text/css; charset=utf-8" {t.Fatalf("CSS content type %q",got)}
 if got:=w.Body.String();got!=suiteButtonCSS {t.Fatal("served CSS must exactly match Homepage embedded module")}
 if len(suiteButtonCSS)==0 {t.Fatal("canonical CSS is empty")}
 if strings.Contains(suiteButtonCSS,"@import") {t.Fatal("runtime CSS import is not permitted")}
 wrong:=httptest.NewRecorder()
 s.ServeHTTP(wrong,httptest.NewRequest(http.MethodPost,"/suite-ui/buttons.css",nil))
 if wrong.Code!=http.StatusMethodNotAllowed {t.Fatalf("POST to CSS returned %d",wrong.Code)}
}

func TestCanonicalButtonStyleAllowedByCSP(t *testing.T) {
 s:=&State{}
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,httptest.NewRequest(http.MethodGet,"/",nil))
 if w.Code!=http.StatusOK {t.Fatalf("page status %d",w.Code)}
 if !strings.Contains(w.Header().Get("Content-Security-Policy"),"style-src 'self' 'unsafe-inline'") {t.Fatal("NewsGo CSP must permit its own embedded stylesheet")}
 if !strings.Contains(w.Body.String(),`href="/suite-ui/buttons.css"`) {t.Fatal("missing canonical stylesheet link")}
}
