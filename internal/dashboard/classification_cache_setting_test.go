package dashboard

import (
 "encoding/json"
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestClassificationCacheSettingPersistenceAndMigration(t *testing.T) {
 def:=defaultArticleSettings()
 if def.ClassificationCacheMinutes!=30 {t.Fatalf("default %d",def.ClassificationCacheMinutes)}
 for _,n:=range []int{0,1,30}{
  v:=def;v.ClassificationCacheMinutes=n
  if !validArticleSettings(v){t.Fatalf("valid value %d rejected",n)}
 }
 for _,n:=range []int{-1,31}{
  v:=def;v.ClassificationCacheMinutes=n
  if validArticleSettings(v){t.Fatalf("out of bounds %d accepted",n)}
 }
 dir:=t.TempDir()
 store,err:=NewSyncStore(dir);if err!=nil {t.Fatal(err)}
 legacy:=defaultServerSettings()
 legacy.Articles.ClassificationCacheMinutes=0
 raw,err:=json.Marshal(legacy);if err!=nil {t.Fatal(err)}
 var document map[string]any
 if err=json.Unmarshal(raw,&document);err!=nil {t.Fatal(err)}
 articles:=document["articles"].(map[string]any)
 delete(articles,"classificationCacheMinutes")
 raw,err=json.Marshal(document);if err!=nil {t.Fatal(err)}
 if err=os.WriteFile(filepath.Join(dir,"server-settings.json"),raw,0600);err!=nil {t.Fatal(err)}
 migrated,err:=store.loadServerSettings();if err!=nil{t.Fatal(err)}
 if migrated.Articles.ClassificationCacheMinutes!=30 {t.Fatalf("legacy migration %d",migrated.Articles.ClassificationCacheMinutes)}
 explicit:=defaultServerSettings();explicit.Articles.ClassificationCacheMinutes=0
 raw,err=json.Marshal(explicit);if err!=nil{t.Fatal(err)}
 if err=os.WriteFile(filepath.Join(dir,"server-settings.json"),raw,0600);err!=nil{t.Fatal(err)}
 reloaded,err:=store.loadServerSettings();if err!=nil{t.Fatal(err)}
 if reloaded.Articles.ClassificationCacheMinutes!=0 {t.Fatalf("explicit zero lost: %d",reloaded.Articles.ClassificationCacheMinutes)}
}

func TestClassificationCacheAdminPresentation(t *testing.T) {
 for _,fragment:=range []string{"admin-classification-cache","classificationCacheMinutes:'admin-classification-cache'","min='0'","max='30'"}{
  if !strings.Contains(personalJS,fragment){t.Errorf("missing control: %s",fragment)}
 }
}
