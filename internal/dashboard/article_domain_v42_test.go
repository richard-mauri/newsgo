package dashboard

import (
	"encoding/json"
	"net/http"
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
		"Exclude articles from ",
		"articleDomains=[...new Set(data.excludedArticleDomains||[])]",
	} {
		if !strings.Contains(personalJS, want) {
			t.Errorf("missing article-domain logic %q", want)
		}
	}
	for _, want := range []string{`id="article-exclude-domain"`, `id="article-exclude-add"`, `Excluded Domains`} {
		if !strings.Contains(page.Tree.Root.String(), want) {
			t.Errorf("missing control %q", want)
		}
	}
}

func TestArticleDomainExclusionsPersistWithUserProfile(t *testing.T) {
	store, e := NewSyncStore(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	token, e := store.IssueToken("richard", "")
	if e != nil {
		t.Fatal(e)
	}
	in := storedPrefs{ExcludedDomains: []string{"nytimes.com"}, ExcludedArticleDomains: []string{"youtube.com"}}
	w := syncCall(t, store, http.MethodPut, token, in)
	if w.Code != 200 {
		t.Fatalf("save: %v %s", w.Code, w.Body.String())
	}
	w = syncCall(t, store, http.MethodGet, token, nil)
	if w.Code != 200 {
		t.Fatalf("load %v", w.Code)
	}
	var out storedPrefs
	if e := json.Unmarshal(w.Body.Bytes(), &out); e != nil {
		t.Fatal(e)
	}
	if out.Revision != 1 || len(out.ExcludedDomains) != 1 || len(out.ExcludedArticleDomains) != 1 || out.ExcludedArticleDomains[0] != "youtube.com" {
		t.Fatalf("persistence: %#v", out)
	}
	in.Revision = 1
	in.ExcludedArticleDomains = []string{"youtube.com", "example.com"}
	w = syncCall(t, store, http.MethodPut, token, in)
	if w.Code != 200 {
		t.Fatalf("revision update: %v %s", w.Code, w.Body.String())
	}
	w = syncCall(t, store, http.MethodPut, token, in)
	if w.Code != 409 {
		t.Fatalf("stale save: %v", w.Code)
	}
}
