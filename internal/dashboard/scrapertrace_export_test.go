package dashboard

import (
 "strings"
 "testing"
)

func TestScraperTraceURLExportWiredToCompletedReadAll(t *testing.T) {
 for _,needle:=range []string{
  "let scraperTraceBatch=[];",
  "if(activeReadingSelection==='all'&&!leadOnlyAllSubscriptions) scraperTraceBatch=eligibleScraperTraceURLs(allPersonalArticles)",
  "el('download-scrapergo-urls').hidden=scraperTraceBatch.length===0",
  "el('download-scrapergo-urls').addEventListener('click',downloadScraperTraceBatch)",
  "if(!article||typeof article.url!=='string'||excludedURL(article.url)||articleExcluded(article.url))continue;",
  "u.hash='';",
  "seen.add(u.href);out.push(u.href);",
  "schema:'newsgo.scrapergo-batch.v1'",
 } {
  if !strings.Contains(personalJS,needle) {t.Errorf("missing diagnostic export contract %q",needle)}
 }
 if !strings.Contains(page.Tree.Root.String(),"id=\"download-scrapergo-urls\" hidden") {
  t.Fatal("missing opt-in download button in briefing")
 }
}
