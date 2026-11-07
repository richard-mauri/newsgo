package dashboard

import (
 "encoding/json"
 "errors"
 "os"
 "path/filepath"
 "time"
)

const classificationCacheCapacity = 2048
const classificationCacheFile = "classification-cache.json"
const classificationCacheMaxBytes = 16 * 1024 * 1024

// Disk entries contain only normalized public URL keys, classification and absolute expiry.
// Never cache errors, unknown results, or unavailable upstream responses.
type persistedClassificationEntry struct {
 Value string `json:"value"`
 ExpiresAt time.Time `json:"expiresAt"`
}
type persistedClassificationCache struct {
 Version int `json:"version"`
 Entries map[string]persistedClassificationEntry `json:"entries"`
}
func validCachedClassification(value string) bool {
 return value=="article" || value=="homepage" || value=="unreadable"
}
func classificationTTLMinutes() int {
 n:=currentArticleSettings().ClassificationCacheMinutes
 if n>720 {return 720}
 return n
}
func validClassificationKey(key string) bool {
 if len(key)==0 || len(key)>4096 {return false}
 // Keys are the ScraperGo reader identity followed by a normalized article URL.
 for i:=0;i<len(key);i++ {
  if key[i]=='|' {
   return i>0 && i<2048 && safeURL(key[:i]) && safeURL(key[i+1:]) &&
    normalizeClassificationURL(key[i+1:])==key[i+1:]
  }
 }
 return false
}
// Called only during server startup, after settings have been activated.
func (s *SettingsStore) ActivateClassificationCache() {
 newsClassificationCache.Lock()
 defer newsClassificationCache.Unlock()
 newsClassificationCache.entries=make(map[string]newsClassificationEntry)
 newsClassificationCache.path=s.file(classificationCacheFile)
 raw,err:=os.ReadFile(newsClassificationCache.path)
 if err!=nil {return} // missing/unreadable cache is a cold cache, not an outage
 if len(raw)>classificationCacheMaxBytes {return}
 var disk persistedClassificationCache
 if json.Unmarshal(raw,&disk)!=nil || disk.Version!=1 || len(disk.Entries)>classificationCacheCapacity {return}
 now:=time.Now()
 needsCleanup:=false
 for key,entry:=range disk.Entries {
  if !validClassificationKey(key)||!validCachedClassification(entry.Value) ||
   !now.Before(entry.ExpiresAt) || entry.ExpiresAt.After(now.Add(12*time.Hour)) {
   needsCleanup=true
   continue
  }
  newsClassificationCache.entries[key]=newsClassificationEntry{entry.Value,entry.ExpiresAt}
 }
 // Compact only when entries were rejected. An I/O error leaves valid
 // classifications in memory; no unrelated settings or profile files change.
 if needsCleanup { _=persistClassificationCacheLocked() }
}
// Caller holds newsClassificationCache.Mutex. On failure, in-memory classifications
// remain valid; disk failures never change article eligibility.
func persistClassificationCacheLocked() error {
 if newsClassificationCache.path=="" {return nil}
 disk:=persistedClassificationCache{Version:1,Entries:make(map[string]persistedClassificationEntry)}
 now:=time.Now()
 for key,entry:=range newsClassificationCache.entries {
  if !now.Before(entry.expiry) {
   delete(newsClassificationCache.entries,key)
   continue
  }
  if validClassificationKey(key) && validCachedClassification(entry.value) &&
   !entry.expiry.After(now.Add(12*time.Hour)) {
   disk.Entries[key]=persistedClassificationEntry{entry.value,entry.expiry}
  }
 }
 if len(disk.Entries)>classificationCacheCapacity {return errors.New("classification cache capacity exceeded")}
 if err:=os.MkdirAll(filepath.Dir(newsClassificationCache.path),0700);err!=nil{return err}
 return atomicJSON(newsClassificationCache.path,disk)
}
