package dashboard

import (
 "context"
 "io"
 "net/http"
 "time"
)

// ServeHomepageSession performs the session verification server-to-server.
// The browser only makes a same-origin request to NewsGo; no account token
// or user identity is accepted in query parameters or request bodies.
func (s *State) ServeHomepageSession(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-store")
 w.Header().Set("X-Content-Type-Options", "nosniff")
 if r.Method != http.MethodGet {
  w.Header().Set("Allow","GET")
  http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
  return
 }
 cookie,err:=r.Cookie("mauri_suite_session")
 if err!=nil || cookie.Value=="" {http.Error(w,"session unavailable",http.StatusUnauthorized);return}
 ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second)
 defer cancel()
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,"https://home.themauris.org/api/session/profile",nil)
 if err!=nil {http.Error(w,"session service unavailable",http.StatusBadGateway);return}
 req.AddCookie(&http.Cookie{Name:"mauri_suite_session",Value:cookie.Value})
 client:=&http.Client{Timeout:5*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 response,err:=client.Do(req)
 if err!=nil {http.Error(w,"session service unavailable",http.StatusBadGateway);return}
 defer response.Body.Close()
 if response.StatusCode==http.StatusUnauthorized {http.Error(w,"session invalid",http.StatusUnauthorized);return}
 if response.StatusCode!=http.StatusOK {http.Error(w,"session service unavailable",http.StatusBadGateway);return}
 w.Header().Set("Content-Type","application/json")
 _,_=io.Copy(w,io.LimitReader(response.Body,4096))
}
