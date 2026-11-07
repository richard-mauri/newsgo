package dashboard

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestClassificationCacheSettingPersistence(t *testing.T){
 def:=defaultArticleSettings()
 if def.ClassificationCacheMinutes!=1440{t.Fatalf("default %d",def.ClassificationCacheMinutes)}
 for _,n:=range []int{0,1,30,1440}{
  cfg:=def;cfg.ClassificationCacheMinutes=n
  if !validArticleSettings(cfg){t.Fatalf("valid value %d rejected",n)}
 }
 for _,n:=range []int{-1,1441}{
  cfg:=def;cfg.ClassificationCacheMinutes=n
  if validArticleSettings(cfg){t.Fatalf("out of range %d accepted",n)}
 }
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 fresh,err:=store.loadServerSettings();if err!=nil{t.Fatal(err)}
 if fresh.Articles.ClassificationCacheMinutes!=1440{t.Fatalf("fresh default=%d",fresh.Articles.ClassificationCacheMinutes)}
 cfg:=defaultServerSettings();cfg.Articles.ClassificationCacheMinutes=0
 raw,err:=json.Marshal(cfg);if err!=nil{t.Fatal(err)}
 if err=os.WriteFile(filepath.Join(dir,"server-settings.json"),raw,0600);err!=nil{t.Fatal(err)}
 reloaded,err:=store.loadServerSettings();if err!=nil{t.Fatal(err)}
 if reloaded.Articles.ClassificationCacheMinutes!=0{t.Fatalf("persisted disabled TTL changed: %d",reloaded.Articles.ClassificationCacheMinutes)}
}

func TestClassificationCacheAdminPresentation(t *testing.T){
 for _,fragment:=range []string{"admin-classification-cache","classificationCacheMinutes:'admin-classification-cache'","min='0'","max='1440'"}{
  if !strings.Contains(personalJS,fragment){t.Errorf("missing cache control %q",fragment)}
 }
}
