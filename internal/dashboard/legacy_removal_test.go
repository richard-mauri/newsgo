package dashboard

import (
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestLegacySyncRouteRetired(t *testing.T){
 state:=&State{}
 for _,method:=range []string{http.MethodGet,http.MethodPut,http.MethodPost}{
  req:=httptest.NewRequest(method,"/api/sync",strings.NewReader("{}"))
  req.Header.Set("Authorization","Bearer "+strings.Repeat("a",64))
  rec:=httptest.NewRecorder()
  state.ServeHTTP(rec,req)
  if rec.Code!=http.StatusNotFound{t.Fatalf("legacy route %s returned %d, want 404",method,rec.Code)}
 }
}
func TestSettingsStoreStartsWithoutLegacyFiles(t *testing.T){
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 if err=store.ActivateServerSettings();err!=nil{t.Fatal(err)}
 for _,name:=range []string{"approved-tokens.json","user-preferences.json"}{
  if _,err=os.Stat(filepath.Join(dir,name));!os.IsNotExist(err){t.Fatalf("unwanted legacy data file %q: %v",name,err)}
 }
}
func TestHomepageOnlyFrontendHasNoLegacyTokenFlow(t *testing.T){
 for _,legacy:=range []string{
  "connectSync(", "syncToken=", "syncCredentialKey", "syncConnected", "loadSync(", "saveSync(", "/api/sync",
 }{
  if strings.Contains(personalJS,legacy){t.Errorf("legacy browser login remains: %q",legacy)}
 }
 if strings.Contains(page.Tree.Root.String(),`id="sync-token"`){t.Fatal("legacy token field remains")}
}
