package dashboard

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync/atomic"
 "testing"
 "time"
)

func TestHomepageAdminRoleVerifiedOnServer(t *testing.T) {
 var calls atomic.Int32
 profile:=struct{User,Role,ExpiresAt string}{User:"test-admin",Role:"admin",ExpiresAt:time.Now().Add(time.Hour).UTC().Format(time.RFC3339)}
 home:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls.Add(1)
  if r.URL.Path!="/api/session/profile"{t.Errorf("unexpected profile path %s",r.URL.Path)}
  if r.Header.Get("Authorization")!="Bearer "+strings.Repeat("a",64){w.WriteHeader(http.StatusUnauthorized);return}
  w.Header().Set("Content-Type","application/json")
  _=json.NewEncoder(w).Encode(profile)
 }))
 defer home.Close()
 state:=&State{homepageProfileURL:home.URL+"/api/session/profile"}
 req:=func(token string)*http.Request{
  r:=httptest.NewRequest(http.MethodGet,"/api/admin/settings",nil)
  if token!=""{r.Header.Set("X-NewsGo-Homepage-Authorization","Bearer "+token)}
  return r
 }
 if !state.homepageAdminVerified(req(strings.Repeat("a",64))){t.Fatal("signed-in Homepage admin not verified")}
 profile.Role="user"
 if state.homepageAdminVerified(req(strings.Repeat("a",64))){t.Fatal("normal Homepage user elevated")}
 profile.Role="admin";profile.ExpiresAt=time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
 if state.homepageAdminVerified(req(strings.Repeat("a",64))){t.Fatal("expired Homepage session elevated")}
 profile.ExpiresAt=time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
 if state.homepageAdminVerified(req(strings.Repeat("b",64))){t.Fatal("revoked/unknown Homepage token elevated")}
 before:=calls.Load()
 legacy:=httptest.NewRequest(http.MethodGet,"/api/admin/settings",nil)
 legacy.Header.Set("Authorization","Bearer "+strings.Repeat("c",64))
 if state.homepageAdminVerified(legacy){t.Fatal("legacy token accepted as Homepage token")}
 if calls.Load()!=before{t.Fatal("legacy token was forwarded to Homepage")}
 if state.homepageAdminVerified(req("invalid")){t.Fatal("malformed token elevated")}
 if calls.Load()!=before{t.Fatal("malformed token sent upstream")}
 home.Close()
 if state.homepageAdminVerified(req(strings.Repeat("a",64))){t.Fatal("network failure allowed admin")}
}

func TestHomepageRoleEnforcedForAdminSettingsWithoutLegacyToken(t *testing.T){
 role:="user"
 home:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  _=json.NewEncoder(w).Encode(map[string]string{"user":"reader","role":role,"expiresAt":time.Now().Add(time.Hour).UTC().Format(time.RFC3339)})
 }))
 defer home.Close()
 store,err:=NewSettingsStore(t.TempDir());if err!=nil{t.Fatal(err)}
 state:=&State{settingsStore:store,homepageProfileURL:home.URL}
 request:=func(method string,payload []byte)*httptest.ResponseRecorder{
  req:=httptest.NewRequest(method,"https://newsgo.themauris.org/api/admin/settings",bytes.NewReader(payload))
  req.Header.Set("X-NewsGo-Homepage-Authorization","Bearer "+strings.Repeat("a",64))
  if method==http.MethodPut{req.Header.Set("Origin","https://newsgo.themauris.org")}
  w:=httptest.NewRecorder();state.ServeHTTP(w,req);return w
 }
 if got:=request("GET",nil).Code;got==http.StatusOK{t.Fatal("normal user accessed admin settings")}
 if got:=request("PUT",[]byte("{}")).Code;got==http.StatusOK{t.Fatal("normal user saved admin settings")}
 role="admin"
 if got:=request("GET",nil).Code;got!=http.StatusOK{t.Fatalf("Homepage admin GET status %d",got)}
}
