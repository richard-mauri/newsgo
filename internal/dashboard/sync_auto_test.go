package dashboard
import("strings";"testing")
func TestAutomaticSyncControlsAndConflictProtection(t *testing.T){
 html:=page.Tree.Root.String()
 for _,x:=range []string{`id="sync-advanced"`,`id="sync-connection"`,`id="sync-load"`,`id="sync-save"`,`id="sync-disconnect"`, `/personal.js?v=0.2.44`} { if !strings.Contains(html,x){t.Errorf("missing sync UI %q",x)}}
 for _,x:=range []string{"newsgo.sync.token.v1","localStorage.setItem(syncCredentialKey,token)","localStorage.getItem(syncCredentialKey)","localStorage.removeItem(syncCredentialKey)","if(needsSave&&syncConnected)scheduleServerSave()", "setInterval(checkSync,30000)","visibilitychange", "syncDirty||syncConflict", "response.status===409", "function saveExclusions(list)","function saveSubscriptions(next)"}{if !strings.Contains(personalJS,x){t.Errorf("missing auto-sync contract %q",x)}}
}
