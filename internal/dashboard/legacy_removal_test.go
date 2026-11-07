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
func TestSettingsStoreDoesNotAlterHistoricalUserFiles(t *testing.T){
 dir:=t.TempDir()
 historical:=map[string]string{
  "approved-tokens.json":`[{"user":"retired","digest":"historical"}]`,
  "user-preferences.json":`{"retired":{"revision":5}}`,
 }
 for name,data:=range historical{
  if err:=os.WriteFile(filepath.Join(dir,name),[]byte(data),0600);err!=nil{t.Fatal(err)}
 }
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 if err:=store.ActivateServerSettings();err!=nil{t.Fatal(err)}
 for name,want:=range historical{
  got,err:=os.ReadFile(filepath.Join(dir,name));if err!=nil{t.Fatal(err)}
  if string(got)!=want{t.Fatalf("historical %s modified",name)}
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
