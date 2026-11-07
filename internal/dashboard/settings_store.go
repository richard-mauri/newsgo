package dashboard

import (
 "encoding/json"
 "errors"
 "os"
 "path/filepath"
 "sync"
)

// SettingsStore persists operational server settings independently of identity.
// Existing data/server-settings.json remains in place during the migration.
type SettingsStore struct {
 mu sync.Mutex
 dir string
}

func NewSettingsStore(dir string) (*SettingsStore,error) {
 if dir=="" {return nil,errors.New("settings directory is empty")}
 if err:=os.MkdirAll(dir,0700);err!=nil{return nil,err}
 return &SettingsStore{dir:dir},nil
}

func(s *SettingsStore) file(name string)string{return filepath.Join(s.dir,name)}

// Atomic updates preserve the existing revisioned settings schema and file.
func atomicJSON(path string,v any)error{
 b,err:=json.MarshalIndent(v,"","  ")
 if err!=nil{return err}
 b=append(b,'\n')
 f,err:=os.CreateTemp(filepath.Dir(path),".newsgo-*.tmp")
 if err!=nil{return err}
 defer os.Remove(f.Name())
 if err=f.Chmod(0600);err!=nil{f.Close();return err}
 if _,err=f.Write(b);err!=nil{f.Close();return err}
 if err=f.Sync();err!=nil{f.Close();return err}
 if err=f.Close();err!=nil{return err}
 return os.Rename(f.Name(),path)
}
