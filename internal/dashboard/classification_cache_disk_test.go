package dashboard

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "sync/atomic"
 "testing"
 "time"
)

func TestClassificationCacheDurableRestartExpiryAndBypass(t *testing.T) {
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 oldSettings:=currentArticleSettings()
 cfg:=defaultArticleSettings();cfg.ClassificationCacheMinutes=720
 activeArticleSettings.Store(cfg)
 defer activeArticleSettings.Store(oldSettings)
 newsClassificationCache.Lock()
 oldPath:=newsClassificationCache.path
 oldEntries:=newsClassificationCache.entries
 newsClassificationCache.path=""
 newsClassificationCache.entries=map[string]newsClassificationEntry{}
 newsClassificationCache.Unlock()
 defer func(){
  newsClassificationCache.Lock()
  newsClassificationCache.path=oldPath
  newsClassificationCache.entries=oldEntries
  newsClassificationCache.Unlock()
 }()
 store.ActivateClassificationCache()
 var calls atomic.Int64
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls.Add(1)
  w.Header().Set("Content-Type","application/json")
  _,_=w.Write([]byte(`{"classification":"article"}`))
 }))
 defer upstream.Close()
 url:="https://example.org/article#section"
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,url,true);got!="article"{t.Fatal(got)}
 if calls.Load()!=1 {t.Fatal("first call did not classify")}
 disk,err:=os.ReadFile(filepath.Join(dir,classificationCacheFile));if err!=nil{t.Fatal(err)}
 if len(disk)==0 {t.Fatal("empty disk cache")}
 // A new process starts with empty memory, but an existing cache file.
 store.ActivateClassificationCache()
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,"https://example.org/article#another",true);got!="article"{t.Fatal(got)}
 if calls.Load()!=1 {t.Fatal("restart missed normalized URL cache")}
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,url,false);got!="article"{t.Fatal(got)}
 if calls.Load()!=2 {t.Fatal("reader opt-out did not bypass cache")}
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,url,true);got!="article"{t.Fatal(got)}
 if calls.Load()!=2 {t.Fatal("bypass erased cached decision")}
 // Simulate expiry on disk, then restart. No sleeps required.
 var saved persistedClassificationCache
 if err=json.Unmarshal(disk,&saved);err!=nil{t.Fatal(err)}
 for k,v:=range saved.Entries{v.ExpiresAt=time.Now().Add(-time.Minute);saved.Entries[k]=v}
 if err=atomicJSON(filepath.Join(dir,classificationCacheFile),saved);err!=nil{t.Fatal(err)}
 store.ActivateClassificationCache()
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,url,true);got!="article"{t.Fatal(got)}
 if calls.Load()!=3 {t.Fatal("expired disk entry was reused")}
 if err=os.WriteFile(filepath.Join(dir,classificationCacheFile),[]byte("{invalid"),0600);err!=nil{t.Fatal(err)}
 store.ActivateClassificationCache()
 if got:=newsClassifyWithCache(context.Background(),upstream.URL,url,true);got!="article"{t.Fatal(got)}
 if calls.Load()!=4 {t.Fatal("corrupt cache did not fail closed")}
 if _,err=os.Stat(filepath.Join(dir,"server-settings.json"));!os.IsNotExist(err){t.Fatal("cache wrote or altered server-settings.json")}
}

func TestClassificationCacheRejectsBadEntriesAndBounds(t *testing.T) {
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 newsClassificationCache.Lock()
 oldPath,oldEntries:=newsClassificationCache.path,newsClassificationCache.entries
 newsClassificationCache.path=""
 newsClassificationCache.entries=map[string]newsClassificationEntry{}
 newsClassificationCache.Unlock()
 defer func(){newsClassificationCache.Lock();newsClassificationCache.path=oldPath;newsClassificationCache.entries=oldEntries;newsClassificationCache.Unlock()}()
 now:=time.Now()
 saved:=persistedClassificationCache{Version:1,Entries:map[string]persistedClassificationEntry{
  "https://reader.example|https://example.org/valid":{Value:"article",ExpiresAt:now.Add(6*time.Hour)},
  "https://reader.example|https://example.org/bad":{Value:"unknown",ExpiresAt:now.Add(30*time.Minute)},
  "https://reader.example|https://example.org/future":{Value:"article",ExpiresAt:now.Add(13*time.Hour)},
  "invalid":{Value:"article",ExpiresAt:now.Add(30*time.Minute)},
 }}
 if err=atomicJSON(filepath.Join(dir,classificationCacheFile),saved);err!=nil{t.Fatal(err)}
 store.ActivateClassificationCache()
 newsClassificationCache.Lock()
 n:=len(newsClassificationCache.entries)
 newsClassificationCache.Unlock()
 if n!=1{t.Fatalf("restored %d entries, want only valid entry",n)}
 // A full snapshot is rejected rather than allocating unbounded memory.
 for i:=0;i<classificationCacheCapacity;i++{
  saved.Entries["https://reader.example|https://example.org/"+string(rune(0x1000+i))]=persistedClassificationEntry{Value:"article",ExpiresAt:now.Add(30*time.Minute)}
 }
 if err=atomicJSON(filepath.Join(dir,classificationCacheFile),saved);err!=nil{t.Fatal(err)}
 store.ActivateClassificationCache()
 newsClassificationCache.Lock();n=len(newsClassificationCache.entries);newsClassificationCache.Unlock()
 if n!=0{t.Fatalf("over-capacity snapshot restored %d entries",n)}
}

func TestClassificationCacheEffectiveTTL(t *testing.T) {
 old:=currentArticleSettings()
 defer activeArticleSettings.Store(old)
 cfg:=defaultArticleSettings()
 for _,tc:=range []struct{configured,want int}{{0,0},{15,15},{60,60},{720,720},{1440,720}}{
  cfg.ClassificationCacheMinutes=tc.configured;activeArticleSettings.Store(cfg)
  if got:=classificationTTLMinutes();got!=tc.want{t.Fatalf("configured %d effective %d want %d",tc.configured,got,tc.want)}
 }
}

func TestClassificationCacheDoesNotPersistFailures(t *testing.T) {
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir);if err!=nil{t.Fatal(err)}
 oldSettings:=currentArticleSettings()
 activeArticleSettings.Store(defaultArticleSettings())
 defer activeArticleSettings.Store(oldSettings)
 newsClassificationCache.Lock()
 oldPath,oldEntries:=newsClassificationCache.path,newsClassificationCache.entries
 newsClassificationCache.path=""
 newsClassificationCache.entries=map[string]newsClassificationEntry{}
 newsClassificationCache.Unlock()
 defer func(){newsClassificationCache.Lock();newsClassificationCache.path=oldPath;newsClassificationCache.entries=oldEntries;newsClassificationCache.Unlock()}()
 store.ActivateClassificationCache()
 for _,reply:=range []string{`{"classification":"unknown"}`,`{"classification":"unavailable"}`,`{"classification":"something-else"}`,`bad-json`} {
  upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
   w.Header().Set("Content-Type","application/json")
   _,_=w.Write([]byte(reply))
  }))
  _=newsClassifyWithCache(context.Background(),upstream.URL,"https://example.org/failure",true)
  upstream.Close()
 }
 store.ActivateClassificationCache()
 newsClassificationCache.Lock()
 n:=len(newsClassificationCache.entries)
 newsClassificationCache.Unlock()
 if n!=0 {t.Fatalf("cached %d unknown/unavailable/failure results",n)}
}

func TestClassificationCacheStartupCompactsExpiredEntries(t *testing.T) {
 dir:=t.TempDir()
 store,err:=NewSettingsStore(dir)
 if err!=nil{t.Fatal(err)}
 newsClassificationCache.Lock()
 oldPath,oldEntries:=newsClassificationCache.path,newsClassificationCache.entries
 newsClassificationCache.path=""
 newsClassificationCache.entries=map[string]newsClassificationEntry{}
 newsClassificationCache.Unlock()
 defer func(){
  newsClassificationCache.Lock()
  newsClassificationCache.path=oldPath
  newsClassificationCache.entries=oldEntries
  newsClassificationCache.Unlock()
 }()
 now:=time.Now()
 validKey:="https://reader.example|https://example.org/current"
 expiredKey:="https://reader.example|https://example.org/expired"
 saved:=persistedClassificationCache{Version:1,Entries:map[string]persistedClassificationEntry{
  validKey:{Value:"article",ExpiresAt:now.Add(11*time.Hour)},
  expiredKey:{Value:"homepage",ExpiresAt:now.Add(-time.Minute)},
 }}
 file:=filepath.Join(dir,classificationCacheFile)
 if err=atomicJSON(file,saved);err!=nil{t.Fatal(err)}
 store.ActivateClassificationCache()
 newsClassificationCache.Lock()
 count:=len(newsClassificationCache.entries)
 _,valid:=newsClassificationCache.entries[validKey]
 newsClassificationCache.Unlock()
 if count!=1||!valid{t.Fatalf("memory entries=%d valid=%v",count,valid)}
 raw,err:=os.ReadFile(file)
 if err!=nil{t.Fatal(err)}
 var cleaned persistedClassificationCache
 if err=json.Unmarshal(raw,&cleaned);err!=nil{t.Fatal(err)}
 if len(cleaned.Entries)!=1{t.Fatalf("disk still has %d entries after startup cleanup",len(cleaned.Entries))}
 if _,exists:=cleaned.Entries[expiredKey];exists{t.Fatal("expired entry remained on disk")}
 if _,exists:=cleaned.Entries[validKey];!exists{t.Fatal("valid entry removed from disk")}
 // A second startup preserves the cleaned, still-fresh entry.
 store.ActivateClassificationCache()
 newsClassificationCache.Lock()
 count=len(newsClassificationCache.entries)
 newsClassificationCache.Unlock()
 if count!=1{t.Fatalf("cleaned cache not reusable after second startup: %d",count)}
}
