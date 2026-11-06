package dashboard

import (
 "strings"
 "testing"
)

func TestScraperTraceURLExportWiredToCompletedReadAll(t *testing.T) {
 for _,needle:=range []string{
  "let scraperTraceBatch=[];",
  "if(activeReadingSelection==='all'&&!leadOnlyAllSubscriptions) scraperTraceBatch=eligibleScraperTraceURLs(allPersonalArticles)",
  "updateScraperTraceExportVisibility();",
  "action.hidden=!(homepageIdentity&&homepageIdentity.role==='admin'&&suiteHandoff&&suiteExpiresAt>Date.now()",
  "if(!homepageIdentity||homepageIdentity.role!=='admin'||!suiteHandoff||suiteExpiresAt<=Date.now())",
  "el('download-scrapergo-urls').addEventListener('click',downloadScraperTraceBatch)",
  "if(!article||typeof article.url!=='string'||excludedURL(article.url)||articleExcluded(article.url))continue;",
  "u.hash='';",
  "seen.add(u.href);out.push(u.href);",
  "schema:'newsgo.scrapergo-batch.v1'",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("missing diagnostic export contract %q",needle)}
 }
 if !strings.Contains(page.Tree.Root.String(),"id=\"download-scrapergo-urls\" class=\"article-url-export-button\" hidden>Download Article URLs") {
  t.Fatal("missing opt-in download button in briefing")
 }
 if strings.Contains(page.Tree.Root.String(),">Download ScraperGo URLs</button>") { t.Fatal("old export label remains") }
 if !strings.Contains(page.Tree.Root.String(),".article-url-export-button[hidden]{display:none!important}") {t.Fatal("admin export hidden state not styled")}
 if !strings.Contains(personalJS,"function accountSessionState(){\n updateScraperTraceExportVisibility();") {t.Fatal("signed-in state changes must update export visibility")}
 // The handler must be wired in the unconditional startup section, not
 // inside accountLoad() or any profile/account action.
 listener:="el('download-scrapergo-urls').addEventListener('click',downloadScraperTraceBatch)"
 if strings.Count(personalJS,listener)!=1 {t.Fatal("download listener must be registered exactly once")}
 setup:=strings.Index(personalJS,"window.addEventListener('beforeunload'")
 listenerAt:=strings.Index(personalJS,listener)
 startup:=strings.Index(personalJS,"const rememberedSyncToken=localStorage.getItem(syncCredentialKey)")
 if setup<0||listenerAt<setup||startup<listenerAt {t.Fatal("download listener must be in startup initialization")}
 if !strings.Contains(personalJS,"JSON.stringify(batch,null,2)+'\\n'") {
  t.Fatal("downloaded JSON needs a real newline escape, not a literal backslash-n suffix")
 }

}
