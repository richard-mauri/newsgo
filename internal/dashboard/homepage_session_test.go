package dashboard

import (
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestHomepageSessionRequiresCookieAndGET(t *testing.T) {
 state:=&State{}
 req:=httptest.NewRequest(http.MethodGet,"https://newsgo.themauris.org/api/homepage/session",nil)
 result:=httptest.NewRecorder()
 state.ServeHomepageSession(result,req)
 if result.Code!=http.StatusUnauthorized{t.Fatalf("missing cookie: got %d",result.Code)}
 if result.Header().Get("Cache-Control")!="no-store"{t.Fatal("session response must not be cached")}
 req=httptest.NewRequest(http.MethodPost,"https://newsgo.themauris.org/api/homepage/session",nil)
 req.AddCookie(&http.Cookie{Name:"mauri_suite_session",Value:"not-a-valid-session"})
 result=httptest.NewRecorder()
 state.ServeHomepageSession(result,req)
 if result.Code!=http.StatusMethodNotAllowed{t.Fatalf("wrong method: got %d",result.Code)}
}
