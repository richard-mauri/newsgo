package dashboard

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func adminCall(t *testing.T, store *SettingsStore, method string, authorized bool, payload any) *httptest.ResponseRecorder {
 t.Helper()
 var body []byte
 if payload!=nil {var err error;body,err=json.Marshal(payload);if err!=nil{t.Fatal(err)}}
 req:=httptest.NewRequest(method,"https://newsgo.themauris.org/api/admin/settings",bytes.NewReader(body))
 if method==http.MethodPut{req.Header.Set("Origin","https://newsgo.themauris.org")}
 w:=httptest.NewRecorder()
 store.ServeAdminSettings(w,req,authorized)
 return w
}
func TestAdminSettingsHomepageOnlyAuthorizationAndPersistence(t *testing.T){
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 if got:=adminCall(t,store,"GET",false,nil).Code;got!=403{t.Fatalf("no admin GET=%d",got)}
 if got:=adminCall(t,store,"PUT",false,defaultServerSettings()).Code;got!=403{t.Fatalf("no admin PUT=%d",got)}
 current:=defaultServerSettings();current.Image.MetadataTimeoutMs=4500
 w:=adminCall(t,store,"PUT",true,current)
 if w.Code!=200{t.Fatalf("put=%d %s",w.Code,w.Body.String())}
 if got:=adminCall(t,store,"PUT",true,current).Code;got!=409{t.Fatalf("stale PUT=%d",got)}
 w=adminCall(t,store,"GET",true,nil)
 if w.Code!=200{t.Fatalf("GET=%d",w.Code)}
 var got serverSettings
 if err=json.Unmarshal(w.Body.Bytes(),&got);err!=nil{t.Fatal(err)}
 if got.Revision!=1||got.Image.MetadataTimeoutMs!=4500{t.Fatalf("settings=%#v",got)}
 if _,err=os.Stat(filepath.Join(dir,"server-settings.json"));err!=nil{t.Fatal(err)}
 if _,err=os.Stat(filepath.Join(dir,"approved-tokens.json"));!os.IsNotExist(err){t.Fatalf("settings write must not create token file: %v",err)}
}
func TestAdminSettingsRangeValidation(t *testing.T){
 store,_:=NewSettingsStore(t.TempDir())
 cfg:=defaultServerSettings();cfg.Image.ParallelRequests=100
 if got:=adminCall(t,store,"PUT",true,cfg).Code;got!=400{t.Fatalf("bad settings=%d",got)}
}
func TestClassificationTimeoutPersistsAcrossRestart(t *testing.T){
 previous:=currentArticleSettings();defer activeArticleSettings.Store(previous)
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 cfg:=defaultServerSettings();cfg.Articles.ClassificationTimeoutSeconds=30
 if got:=adminCall(t,store,"PUT",true,cfg).Code;got!=200{t.Fatalf("save=%d",got)}
 restarted,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 if err=restarted.ActivateServerSettings();err!=nil{t.Fatal(err)}
 if got:=currentArticleSettings().ClassificationTimeoutSeconds;got!=30{t.Fatalf("timeout=%d",got)}
 if got:=adminCall(t,restarted,"GET",true,nil).Code;got!=200{t.Fatalf("reload=%d",got)}
}
func TestAdminSettingsResponsiveLayoutContract(t *testing.T){
 html:=page.Tree.Root.String()
 for _,want:=range []string{"admin-settings-grid","admin-settings-actions","admin-feed-search-timeout","grid-template-columns:repeat(2,minmax(0,1fr))","@media(max-width:700px)"}{
  if !strings.Contains(html,want){t.Errorf("missing admin layout %q",want)}
 }
}
func TestUnsavedAdminSettingsDialogContract(t *testing.T){
 html:=page.Tree.Root.String()
 for _,want:=range []string{"admin-unsaved-dialog","Save changes","Discard changes","Keep editing"}{
  if !strings.Contains(html,want){t.Errorf("missing unsaved dialog %q",want)}
 }
 for _,want:=range []string{"const adminDirty=","confirmAdminLeave(","await saveAdmin()","beforeunload"}{
  if !strings.Contains(personalJS,want){t.Errorf("missing dirty protection %q",want)}
 }
}
func TestLegacyExitAndCredentialLoginAbsent(t *testing.T){
 html:=page.Tree.Root.String()
 for _,legacy:=range []string{"admin-exit-dialog","id=\"sync-token\"","Legacy Server Sync"}{
  if strings.Contains(html,legacy){t.Fatalf("legacy UI remains: %q",legacy)}
 }
 for _,legacy:=range []string{"newsgo.sync.token.v1","syncCredentialKey","connectSync(","/api/sync"}{
  if strings.Contains(personalJS,legacy){t.Fatalf("legacy client code remains: %q",legacy)}
 }
}

func TestThumbnailAdminSettingsPersistAndLegacyDefault(t *testing.T) {
 old:=currentThumbnailSettings()
 defer activeThumbnailSettings.Store(old)
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 candidate:=defaultServerSettings()
 candidate.Thumbnail=ThumbnailSettings{Concurrency:4,QueueWaitSeconds:8,UpstreamTimeoutSeconds:15}
 if w:=adminCall(t,store,http.MethodPut,true,candidate);w.Code!=http.StatusOK{t.Fatalf("save=%d %s",w.Code,w.Body.String())}
 restarted,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 if err=restarted.ActivateServerSettings();err!=nil{t.Fatal(err)}
 if got:=currentThumbnailSettings();got!=candidate.Thumbnail{t.Fatalf("thumbnail settings lost: %+v",got)}
 bad:=restarted.CurrentServerSettings()
 bad.Thumbnail.Concurrency=17
 if got:=adminCall(t,restarted,http.MethodPut,true,bad).Code;got!=http.StatusBadRequest{t.Fatalf("invalid concurrency allowed: %d",got)}
 // A previous-version saved document has no thumbnail field.
 legacy:=defaultServerSettings()
 legacy.Thumbnail=ThumbnailSettings{}
 if err=atomicJSON(filepath.Join(dir,"server-settings.json"),legacy);err!=nil{t.Fatal(err)}
 if err=restarted.ActivateServerSettings();err!=nil{t.Fatal(err)}
 if got:=currentThumbnailSettings();got!=defaultThumbnailSettings(){t.Fatalf("legacy settings not defaulted: %+v",got)}
 if _,err=os.Stat(filepath.Join(dir,"classification-cache.json"));!os.IsNotExist(err){t.Fatal("legacy load unexpectedly modified cache")}
}

func TestThumbnailAdminLayoutAndDiagnosticsContract(t *testing.T) {
 html:=page.Tree.Root.String()
 for _,item:=range []string{"admin-thumb-concurrency","admin-thumb-queue","admin-thumb-timeout","thumbnail-performance-panel","thumbnail-performance"}{
  if !strings.Contains(html,item){t.Errorf("missing thumbnail field %s",item)}
 }
 for _,item:=range []string{"thumbnailAdminFields","thumbnailServerTiming(response)","Queue wait mean / p95","Upstream mean / p95","thumbnailCanceledCount"}{
  if !strings.Contains(personalJS,item){t.Errorf("missing thumbnail behavior %s",item)}
 }
}
