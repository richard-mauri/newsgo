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
  "stageReaderSetting(readingSelectionKey,value)",
  "stageReaderSetting(readingSelectionKey,null)",
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
