package dashboard

import (
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 suiteui "github.com/richard-mauri/homepage/shared/ui"
)

var suiteButtonCSS=string(suiteui.ButtonCSS)

func TestCanonicalButtonStylesheetServedLocally(t *testing.T) {
 s:=&State{}
 r:=httptest.NewRequest(http.MethodGet,"/suite-ui/buttons.css",nil)
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,r)
 if w.Code!=http.StatusOK {t.Fatalf("CSS status %d",w.Code)}
 if got:=w.Header().Get("Content-Type");got!="text/css; charset=utf-8" {t.Fatalf("CSS content type %q",got)}
 if got:=w.Body.String();got!=suiteButtonCSS {t.Fatal("served CSS must exactly match Homepage embedded module")}
 if len(suiteButtonCSS)==0 {t.Fatal("canonical CSS is empty")}
 if strings.Contains(suiteButtonCSS,"@import") {t.Fatal("runtime CSS import is not permitted")}
 wrong:=httptest.NewRecorder()
 s.ServeHTTP(wrong,httptest.NewRequest(http.MethodPost,"/suite-ui/buttons.css",nil))
 if wrong.Code!=http.StatusMethodNotAllowed {t.Fatalf("POST to CSS returned %d",wrong.Code)}
}

func TestCanonicalButtonStyleAllowedByCSP(t *testing.T) {
 s:=&State{}
 w:=httptest.NewRecorder()
 s.ServeHTTP(w,httptest.NewRequest(http.MethodGet,"/",nil))
 if w.Code!=http.StatusOK {t.Fatalf("page status %d",w.Code)}
 if !strings.Contains(w.Header().Get("Content-Security-Policy"),"style-src 'self' 'unsafe-inline'") {t.Fatal("NewsGo CSP must permit its own embedded stylesheet")}
 if !strings.Contains(w.Body.String(),`href="/suite-ui/buttons.css"`) {t.Fatal("missing canonical stylesheet link")}
}

func TestNoRedundantLocalButtonAppearance(t *testing.T) {
 markup:=page.Tree.Root.String()
 // Placement rules remain local, but regular-button appearance belongs to
 // Homepage. The subscription-read link is an intentional semantic exception.
 for _,old:=range []string{
  ".catalog-search-form button{", ".catalog-card button{", ".catalog-card button:disabled{",
  ".personal-tools button,.personal-form button,.personal-row button{",
  ".sync-choice-actions button{", ".admin-settings-actions button{",
  ".news-nav button:focus-visible{", ".filtered-domain-exclude:hover{",
  ".filtered-domain-exclude:focus-visible{", ".article-url-export-button:hover{",
  ".article-url-export-button:focus-visible{", ".article-url-export-button{",
  ".included-card button:disabled{", "button.ui-button{appearance:none",
 } {
  if strings.Contains(markup,old) {t.Errorf("obsolete local button styling remains: %s",old)}
 }
 if !strings.Contains(markup,".included-card button{width:100%;min-width:0}") {t.Error("subscription column layout was lost")}
 if !strings.Contains(markup,".discovered-feed-group-heading button.ui-button{width:205px;max-width:100%}") {t.Error("publisher column layout was lost")}
 if !strings.Contains(markup,".subscription-read{display:block;border:0!important") {t.Error("intentional inline-link action style was lost")}
 for _,want:=range []string{"confirm.classList.add('exclusion-confirm-action')", "action.classList.add('subscription-unsubscribe-button')", "exclude.classList.add('filtered-domain-exclude')"} {
  if !strings.Contains(personalJS,want) {t.Errorf("dynamic button lost shared class: %s",want)}
 }
 for _,old:=range []string{"background:#18252d;color:#f1f5f9", "background:transparent;color:var(--muted)"} {
  if strings.Contains(personalJS,old) {t.Errorf("legacy inline button styling remains: %s",old)}
 }
}
