package dashboard

import (
 "os"
 "strings"
 "testing"
)

// A public HTML handoff shell must not show reader controls or make data calls
// until Homepage has verified an actual user or Guest child bearer.
func TestNewsGoBrowserRequiresHomepageIdentity(t *testing.T) {
 for _,want:=range []string{
  "html:not(.suite-auth-ready) .suite-content{visibility:hidden}",
 } {
  if !strings.Contains(page.Tree.Root.String(),want){
   t.Errorf("missing inert bootstrap rule %q",want)
  }
 }
 for _,want:=range []string{
  "function allowVerifiedNewsgo(){document.documentElement.classList.add('suite-auth-ready')}",
  "function denyUnverifiedNewsgo(){document.documentElement.classList.remove('suite-auth-ready')}",
  "rejectedLegacyLocalMode=true",
  "location.replace(homepageOrigin+'/#suite_signin=required')",
  "async function authorizedNewsgoFetch(path,options={})",
  "if(!homepageIdentity||!suiteHandoff||suiteExpiresAt<=Date.now()||rejectedLegacyLocalMode)",
  "Authorization:'Bearer '+suiteHandoff",
  "if(response.status===401||response.status===403)",
  "if(response.status===503)",
  "if(!['guest','user','admin'].includes(profile.role))",
 } {
  if !strings.Contains(personalJS,want){t.Errorf("missing browser authorization contract %q",want)}
 }
 for _,path:=range []string{
  "/api/search-feeds?","/api/personal","/api/discover",
  "/api/admin/settings","/api/admin/classification-diagnostic",
 } {
  if !strings.Contains(personalJS,"authorizedNewsgoFetch('"+path) {
   t.Errorf("%s does not send verified bearer through the guarded fetch",path)
  }
  if strings.Contains(personalJS,"fetch('"+path) && !strings.Contains(personalJS,"authorizedNewsgoFetch('"+path) {
   t.Errorf("%s still uses anonymous fetch",path)
  }
 }
 if strings.Contains(personalJS,"explicitBrowserOnly=true") ||
  strings.Contains(personalJS,"Not signed in · Browser settings only"){
  t.Fatal("legacy unauthenticated reader mode must be rejected")
 }
}

func TestNewsGoExecutableEnablesServerGuard(t *testing.T) {
 source,err:=os.ReadFile("../../cmd/newsgo/main.go")
 if err!=nil{t.Fatal(err)}
 if !strings.Contains(string(source),"dashboard.NewNewsGoAuthGate().Protect(app)"){
  t.Fatal("production server must use the shared guard at its HTTP boundary")
 }
}
