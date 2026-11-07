package dashboard

import (
 "encoding/json"
 "io"
 "net/http"
 "strings"
 "time"
)

const homepageProfileEndpoint = "https://home.themauris.org/api/session/profile"

// homepageAdminVerified checks the bearer against Homepage on every privileged
// request. Never trust the browser's role claim or a cached profile. The
// endpoint is fixed in production; the State override exists for tests only.
func (s *State) homepageAdminVerified(r *http.Request) bool {
 raw:=r.Header.Get("Authorization")
 if !strings.HasPrefix(raw,"Bearer "){return false}
 token:=strings.TrimPrefix(raw,"Bearer ")
 if len(token)!=64{return false}
 for _,c:=range token{if !(c>='0'&&c<='9'||c>='a'&&c<='f'){return false}}
 endpoint:=s.homepageProfileURL
 if endpoint==""{endpoint=homepageProfileEndpoint}
 ctx:=r.Context()
 request,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil)
 if err!=nil{return false}
 request.Header.Set("Authorization","Bearer "+token)
 request.Header.Set("Accept","application/json")
 client:=&http.Client{Timeout:3*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 response,err:=client.Do(request)
 if err!=nil{return false}
 defer response.Body.Close()
 if response.StatusCode!=http.StatusOK{return false}
 var profile struct {
  User string `json:"user"`
  Role string `json:"role"`
  ExpiresAt string `json:"expiresAt"`
 }
 decoder:=json.NewDecoder(io.LimitReader(response.Body,2048))
 if decoder.Decode(&profile)!=nil{return false}
 expires,err:=time.Parse(time.RFC3339,profile.ExpiresAt)
 if err!=nil || !time.Now().Before(expires){return false}
 return profile.User!=""&&profile.Role=="admin"
}
