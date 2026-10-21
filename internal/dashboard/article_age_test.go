package dashboard

import (
 "testing"
 "time"
)

func TestArticleAgeValidation(t *testing.T) {
 for _,v:=range []int{0,1,3,7,14,30,60} {if !validArticleAgeDays(v){t.Errorf("rejected supported age %d",v)}}
 for _,v:=range []int{-1,2,4,61,365} {if validArticleAgeDays(v){t.Errorf("accepted unsupported age %d",v)}}
}
func TestArticleAgeCutoff(t *testing.T) {
 now:=time.Date(2026,10,21,12,0,0,0,time.UTC)
 if articleTooOld(time.Time{},7,now){t.Fatal("undated entry excluded")}
 if articleTooOld(now.Add(-30*24*time.Hour),0,now){t.Fatal("all-articles setting filtered entry")}
 if articleTooOld(now.Add(-6*24*time.Hour),7,now){t.Fatal("recent entry excluded")}
 if !articleTooOld(now.Add(-8*24*time.Hour),7,now){t.Fatal("old entry kept")}
}
func TestAgePreferenceSyncValidation(t *testing.T) {
 for _, v := range []int{0,1,3,7,14,30,60} {if !validPrefs(storedPrefs{MaxArticleAgeDays:v}) {t.Errorf("valid sync preference %d rejected",v)}}
 if validPrefs(storedPrefs{MaxArticleAgeDays:999}) {t.Fatal("invalid age accepted")}
}
