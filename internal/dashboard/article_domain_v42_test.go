package dashboard

import (
	"strings"
	"testing"
)

func TestArticleDomainExclusionPipelineV42(t *testing.T) {
	for _, want := range []string{
		"function articleExcluded(value)",
		"articleDomains.includes('youtube.com')&&hostOf(value)==='youtu.be'", 
		"sources.has(a.source)&&!articleExcluded(a.url)", 
		"excludedArticleDomains:articleDomains",
		"function saveArticleExclusions(next)",
		"function toggleArticleDomain(value){toggleDomain(value)}",
		"articleDomains=[...new Set(data.excludedArticleDomains||[])]",
	} {
		if !strings.Contains(personalJS, want) {
			t.Errorf("missing article-domain logic %q", want)
		}
	}
	for _, want := range []string{`id="exclude-domain"`, `id="exclude-add"`, `Excluded Domains`} {
		if !strings.Contains(page.Tree.Root.String(), want) {
			t.Errorf("missing control %q", want)
		}
	}
 for _,legacy:=range []string{`id="article-exclude-domain"`,`id="article-exclude-add"`} {if strings.Contains(page.Tree.Root.String(),legacy){t.Errorf("dead hidden exclusion control remains: %q",legacy)}}
}

func TestArticleDomainExclusionsPersistWithHomepage(t *testing.T){
 for _,needle:=range []string{
  "excludedArticleDomains:articleDomains",
  "async function accountLoad()",
  "async function accountEstablish()",
  "async function accountSave()",
  "revision:accountRevision",
  "articleDomains=[...new Set(data.excludedArticleDomains||[])]",
 }{
  if !strings.Contains(personalJS,needle){t.Errorf("missing Homepage article exclusion persistence %q",needle)}
 }
}
