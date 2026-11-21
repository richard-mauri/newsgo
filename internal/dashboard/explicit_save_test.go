package dashboard

import (
 "strings"
 "testing"
)

func TestReaderSettingsOnlyCommitOnSaveChanges(t *testing.T) {
 for _,want:=range []string{
  "const pendingReaderSettings=new Map()",
  "function stageReaderSetting(key,value)",
  "async function commitReaderSettings()",
  "el('reader-save-changes').addEventListener('click'",
  "stageReaderSetting(unifiedKey,JSON.stringify(next))",
  "stageReaderSetting(exclusionKey,JSON.stringify(excludedDomains))",
  "stageReaderSetting(articleExclusionKey,JSON.stringify(articleDomains))",
  "stageReaderSetting(agePreferenceKey,String(maxArticleAgeDays))",
  "function rememberReadingSelection(value)",
  "localStorage.setItem(readingSelectionKey,String(value))",
  "rememberReadingSelection(null)",
  "localStorage.removeItem(readingSelectionKey)",
  "stageReaderSetting(leadCountKey,String(maximumLeadStories))",
  "stageReaderSetting(leadOnlyKey,String(leadOnlyAllSubscriptions))",
  "stageReaderSetting(classificationCachePreferenceKey,String(useClassificationCache))",
  "stageReaderSetting(rssOnlyPreferenceKey,String(rssMetadataOnly))",
  "if(adminDirty()||hasUnsavedReaderSettings())",
  "accountSyncNeeded=true",
  "async function accountSave(snapshot)",
  "No saved record; awaiting Save Changes",
  "readerUnsyncedKey(owner)",
  "if(unsyncedSnapshot&&!staged)applyAccountPreferences(unsyncedSnapshot)",
 } {
  if !strings.Contains(personalJS,want){t.Errorf("missing explicit-save contract %q",want)}
 }
 if strings.Contains(personalJS,"stageReaderSetting(readingSelectionKey") {t.Fatal("reading navigation must not enter pending Save Changes settings")} 
 for _,bad:=range []string{
  "scheduleAccountSave()",
  "setTimeout(accountSave,650)",
  "setTimeout(accountSave,30000)",
  "accountDirty=true",
  "if(!record)scheduleAccountSave()",
 }{
  if strings.Contains(personalJS,bad){t.Errorf("automatic persistence reappeared: %q",bad)}
 }
 html:=page.Tree.Root.String()
 for _,want:=range []string{"id=\"reader-save-changes\"","id=\"reader-save-status\"","Save Changes"} {
  if !strings.Contains(html,want){t.Errorf("missing Save Changes interface %q",want)}
 }
 if !strings.Contains(personalJS,"localStorage.setItem(field.key") || !strings.Contains(personalJS,"localStorage.removeItem(field.key)"){
  t.Fatal("browser-only recent-search histories should remain independent of Save Changes")
 }
 if !strings.Contains(personalJS,"async function saveAdmin()") || !strings.Contains(personalJS,"adminDirty()"){
  t.Fatal("administrator server-settings save must remain independent")
 }
}

func TestReaderSaveSharesMainNavigationRow(t *testing.T) {
 markup:=page.Tree.Root.String()
 navStart:=strings.Index(markup,`<nav class="news-nav"`)
 if navStart<0 {t.Fatal("NewsGo navigation missing")}
 navEnd:=strings.Index(markup[navStart:],"</nav>")
 if navEnd<0 {t.Fatal("NewsGo navigation closing tag missing")}
 nav:=markup[navStart:navStart+navEnd]
 for _,id:=range []string{`id="manage-personal"`,`id="view-subscriptions"`,`id="manage-settings"`,`id="manage-diagnostics"`,`id="reader-save-changes"`,`id="reader-save-status"`} {
  if !strings.Contains(nav,id){t.Errorf("navigation row missing %s",id)}
 }
 if strings.Contains(markup,`class="reader-save-bar"`){t.Fatal("Save Changes must not render in a separate toolbar row")}
 if !strings.Contains(markup,`.news-nav{display:flex;align-items:center;gap:10px;flex-wrap:wrap;`){
  t.Fatal("navigation alignment and narrow-screen wrapping missing")
 }
}

func TestHomepageAuthAndExplicitSaveRemainIntegrated(t *testing.T) {
 // Both owner-accepted development scopes must remain in one executable.
 for _,want:=range []string{
  "function allowVerifiedNewsgo()",
  "function denyUnverifiedNewsgo()",
  "async function authorizedNewsgoFetch(path,options={})",
  "Authorization:'Bearer '+suiteHandoff",
  "location.replace(homepageOrigin+'/#suite_signin=required')",
  "if(!homepageIdentity||!suiteHandoff||suiteExpiresAt<=Date.now()||rejectedLegacyLocalMode)",
  "async function commitReaderSettings()",
  "async function accountSave(snapshot)",
  "if(personal)writes.push([readerUnsyncedKey(owner),JSON.stringify(snapshot)])",
  "if(personal){",
  "if(adminDirty()||hasUnsavedReaderSettings())",
 } {
  if !strings.Contains(personalJS,want){t.Errorf("combined auth + Save Changes contract missing %q",want)}
 }
 for _,path:=range []string{"/api/personal","/api/discover","/api/search-feeds?","/api/admin/settings"} {
  if !strings.Contains(personalJS,"authorizedNewsgoFetch('"+path){
   t.Errorf("combined candidate has unguarded application fetch %s",path)
  }
 }
 html:=page.Tree.Root.String()
 if !strings.Contains(html,"html:not(.suite-auth-ready) .suite-content{visibility:hidden}") {
  t.Fatal("unauthenticated reader must remain hidden during handoff")
 }
}
